/*
Package text holds models which
make text classification easy. They
are regular models, but take strings
as arguments so you can feed in
documents rather than large,
hand-constructed word vectors. Although
models might represent the words as
these vectors, the munging of a
document is hidden from the user.

The simplese model, although suprisingly
effective, is Naive Bayes. If you
want to read more about the specific
model, check out the docs for the
NaiveBayes struct/model.

Example Online Naive Bayes Text Classifier (multiclass):

	// create the channel of data and errors
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	// make a new NaiveBayes model with
	// 2 classes expected (classes in
	// datapoints will now expect {0,1}.
	// in general, given n as the classes
	// variable, the model will expect
	// datapoint classes in {0,...,n-1})
	model := NewNaiveBayes(stream, 2)

	go model.OnlineLearn(errors)

	stream <- base.TextDatapoint{
		X: "I love the city"
		Y: 1
	}

	stream <- base.TextDatapoint{
		X: "I hate Los Angeles"
		Y: 0
	}

	stream <- base.TextDatapoint{
		X: "My mother is not a nice lady"
		Y: 0
	}

	close(stream)

	for {
		err, more := <- errors
		if err != nil {
			fmt.Printf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}

	// now you can predict like normal
	class := model.Predict("My mother is in Los Angeles") // 0
*/
package text

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
NaiveBayes is a general classification
model that calculates the probability
that a datapoint is part of a class
by using Bayes Rule:
	P(y|x) = P(x|y)*P(y)/P(x)
The unique part of this model is that
it assumes words are unrelated to
eachother. For example, the probability
of seeing the word 'penis' in spam emails
if you've already seen 'viagra' might be
different than if you hadn't seen it. The
model ignores this fact because the
computation of full Bayesian model would
take much longer, and would grow significantly
with each word you see.

https://en.wikipedia.org/wiki/Naive_Bayes_classifier
http://cs229.stanford.edu/notes/cs229-notes2.pdf

Based on Bayes Rule, we can easily calculate
the numerator (x | y is just the number of
times x is seen and the class=y, and P(y) is
just the number of times y=class / the number
of positive training examples/words.) The
denominator is also easy to calculate, but
if you recognize that it's just a constant
because it's just the probability of seeing
a certain document given the dataset we can
make the following transformation to be
able to classify without as much classification:
	Class(x) = argmax_c{P(y = c) * ∏P(x|y = c)}
And we can use logarithmic transformations to
make this calculation more computer-practical
(multiplying a bunch of probabilities on [0,1]
will always result in a very small number
which could easily underflow the float value):
	Class(x) = argmax_c{log(P(y = c)) + ΣP(x|y = c)}
Much better. That's our model!
*/
type NaiveBayes struct {
	// Words holds a map of words
	// to their corresponding Word
	// structure
	Words map[string]Word `json:"words"`

	// Count holds the number of times
	// class i was seen as Count[i]
	Count []uint64 `json:"count"`

	// Probabilities holds the probability
	// that class Y is class i as
	// Probabilities[i] for
	Probabilities []float64 `json:"probabilities"`

	// DictSize holds the number of
	// words being tracked in the
	// dictionary
	DictSize uint64 `json:"dict_size"`

	// sanitize is used by a model
	// to sanitize input of text
	sanitize func(r rune) bool

	// stream holds the datastream
	stream chan base.TextDatapoint
}

// Word holds the structural
// information needed to calculate
// the probability of
type Word struct {
	// Count holds the number of times,
	// (i in Count[i] is the given class)
	Count []uint64

	// Seen holds the number of times
	// the world has been seen. This
	// is than same as
	//    foldl (+) 0 Count
	// in Haskell syntax, but is included
	// you wouldn't have to calculate
	// this every time you wanted to
	// recalc the probabilities (foldl
	// is the same as reduce, basically.)
	Seen uint64

	// Probabilities holds the probability
	// that word x is of class i as
	// Probabilities[i]
	Probabilities []float64
}

// NewNaiveBayes returns a NaiveBayes model the
// given number of classes instantiated, ready
// to learn off the given data stream. The sanitization
// function is set to the given function. It must
// comply with the transform.RemoveFunc interface
func NewNaiveBayes(stream base.TextDatapoint, classes uint8, sanitize func(rune) bool) *NaiveBayes {
	return &NaiveBayes{
		Words:         make(map[string]Word),
		Count:         make([]uint64, classes),
		Probabilities: make([]float64, classes),
		DictSize:      uint64(0),

		sanitize: sanitize,
		stream:   stream,
	}
}

// UpdateStream updates the NaiveBayes model's
// text datastream
func (b *NaiveBayes) UpdateStream(stream chan base.TextDatapoint) {
	b.stream = stream
}

// UpdateSanitize updates the NaiveBayes model's
// text sanitization transformation function
func (b *NaiveBayes) UpdateSanitize(sanitize func(rune) bool) {
	b.sanitize = sanitize
}

// String implements the fmt interface for clean printing. Here
// we're using it to print the model as the equation h(θ)=...
// where h is the perceptron hypothesis model.
func (b *NaiveBayes) String() string {
	return fmt.Sprintf("h(θ) = argmax_c{log(P(y = c)) + ΣP(x|y = c)}\n\tClasses: %v\n\tWords evaluated in model: %v\n", len(b.Classes), int(b.DictSize))
}

// PersistToFile takes in an absolute filepath and saves the
// parameter vector θ to the file, which can be restored later.
// The function will take paths from the current directory, but
// functions
//
// The data is stored as JSON because it's one of the most
// efficient storage method (you only need one comma extra
// per feature + two brackets, total!) And it's extendable.
func (b *NaiveBayes) PersistToFile(path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to persist your model to a file with no path!! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := json.Marshal(b)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// RestoreFromFile takes in a path to a parameter vector theta
// and assigns the model it's operating on's parameter vector
// to that.
//
// The path must ba an absolute path or a path from the current
// directory
//
// This would be useful in persisting data between running
// a model on data.
func (b *NaiveBayes) RestoreFromFile(path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to restore your model from a file with no path! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &b)
	if err != nil {
		return err
	}

	return nil
}