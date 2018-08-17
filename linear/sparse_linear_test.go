package linear

import (
	//"fmt"
	//"math/rand"
	//"os"
	"testing"

	"github.com/bountylabs/goml/base"

	"github.com/stretchr/testify/assert"
	//"fmt"
	"fmt"
)

var sparseFlatX []map[int]float64
var sparseIncreasingX []map[int]float64
var sparseThreeDLineX []map[int]float64
//var sparseNormX []map[int]float64
var sparseNoisyX []map[int]float64

func init() {

	/*
		var flatX [][]float64
		var flatY []float64

		var increasingX [][]float64
		var increasingY []float64

		var threeDLineX [][]float64
		var threeDLineY []float64

		var normX [][]float64
		var normY []float64

		var noisyX [][]float64
		var noisyY []float64
	*/

	sparseFlatX = denseToSparse(flatX)
	sparseIncreasingX = denseToSparse(increasingX)
	sparseThreeDLineX = denseToSparse(threeDLineX)
	sparseNoisyX = denseToSparse(noisyX)

}

func denseToSparse(in [][]float64) []map[int]float64 {
	out := make([]map[int]float64, len(in))
	for index, features := range in {
		mp := map[int]float64{}
		for findex, v := range features {
			mp[findex] = v
		}
		out[index] = mp
	}
	return out
}

// test y=3
func TestSparseFlatLineShouldPass1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .01, 0.1, 0, base.L2, 800, sparseFlatX, flatY, len(sparseFlatX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				assert.InDelta(t, 3, guess[0], 1e-2, "Guess should be really close to 3 (within 1e-2) for y=3")
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}
}

// same as above but with StochasticGD
func TestSparseFlatLineShouldPass2(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .00001, .00001, 0, base.L2, 800, sparseFlatX, flatY, len(sparseFlatX[0]))

	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				assert.InDelta(t, 3, guess[0], 1e-2, "Guess should be really close to 3 (within 1e-2) for y=3")
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}
}

// test y=3 but don't have enough iterations
func TestSparseFlatLineShouldFail1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .01, 0.1, 0, base.L2, 1, sparseFlatX, flatY, len(sparseFlatX[0]))

	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64
	var faliures int

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				if abs(3.0-guess[0]) > 1e-2 {
					faliures++
				}
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}

	assert.True(t, faliures > 40, "There should be more faliures than half of the training set")
}

// same as above but with StochasticGD
func TestSparseFlatLineShouldFail2(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .001, 0.01, 0, base.L2, 1, sparseFlatX, flatY, len(sparseFlatX[0]))

	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64
	var faliures int

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				if abs(3.0-guess[0]) > 1e-2 {
					faliures++
				}
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}

	assert.True(t, faliures > 40, "There should be more faliures than half of the training set")
}

// test y=3 but include an invalid data set
func TestSparseFlatLineShouldFail3(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .01, 0.1, 0, base.L2, 1, []map[int]float64{}, flatY, len(sparseFlatX[0]))

	err = model.Learn("")
	assert.NotNil(t, err, "Learning error should not be nil")
}

// same as above but with StochasticGD
func TestSparseFlatLineShouldFail4(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .01, 0.1, 0, base.L2, 1, []map[int]float64{}, flatY, len(sparseFlatX[0]))

	err = model.Learn("")
	assert.NotNil(t, err, "Learning error should not be nil")
}

// test y=3 but include an invalid data set
func TestSparseFlatLineShouldFail5(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .01, 0.1, 0, base.L2, 1, []map[int]float64{map[int]float64{}}, nil, len(sparseFlatX[0]))
	err = model.Learn("")
	assert.NotNil(t, err, "Learning error should not be nil")
}

// invalid optimization method
func TestSparseFlatLineShouldFail6(t *testing.T) {
	var err error

	model := NewSparseLeastSquares("not a method", .01, 0.1, 0, base.L2, 1, sparseFlatX, flatY, len(sparseFlatX[0]))
	err = model.Learn("")
	assert.NotNil(t, err, "Learning error should not be nil")
}

// test y=x
func TestSparseInclinedLineShouldPass1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .01, 0.1, 0, base.L2, 500, sparseIncreasingX, increasingY, len(increasingX[0]))

	//model := NewLeastSquares(base.BatchGD, .01, 0, 500, increasingX, increasingY)
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i, guess[0], 1e-2, "Guess should be really close to input (within 1e-2) for y=x")
		assert.Nil(t, err, "Prediction error should be nil")
	}
}

// same as above but with StochasticGD
func TestSparseInclinedLineShouldPass2(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .01, .01, 0, base.L2, 500, sparseIncreasingX, increasingY, len(increasingX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i, guess[0], 1e-2, "Guess should be really close to input (within 1e-2) for y=x")
		assert.Nil(t, err, "Prediction error should be nil")
	}
}

// test y=x but regularization term too large
func TestSparseInclinedLineShouldFail1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .0001, .0001, 1e3, base.L2, 500, sparseIncreasingX, increasingY, len(increasingX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64
	var faliures int

	for i := -20; i < 20; i += 2 {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		if abs(float64(i)-guess[0]) > 1e-2 {
			faliures++
		}
		assert.Nil(t, err, "Prediction error should be nil")
	}

	assert.True(t, faliures > 15, "There should be more faliures than half of the training set")
}

// same as above but with StochasticGD
func TestSparseInclinedLineShouldFail2(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 1e3, base.L2, 500, sparseIncreasingX, increasingY, len(increasingX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64
	var faliures int

	for i := -20; i < 20; i += 2 {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		if abs(float64(i)-guess[0]) > 1e-2 {
			faliures++
		}
		assert.Nil(t, err, "Prediction error should be nil")
	}

	assert.True(t, faliures > 15, "There should be more faliures than half of the training set")
}

// test z = 10 + (x/10) + (y/5)
func TestSparseThreeDimensionalLineShouldPass1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, .01, .01, 0, base.L2, 500, sparseThreeDLineX, threeDLineY, len(threeDLineX[0]))
	err = model.Learn("")

	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			guess, err = model.Predict([]float64{float64(i), float64(j)})
			assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
			assert.InDelta(t, 10.0+float64(i)/10+float64(j)/5, guess[0], 1e-2, "Guess should be really close to i+x (within 1e-2) for line z=10 + (x+y)/10")
			assert.Nil(t, err, "Prediction error should be nil")
		}
	}
}

// same as above but with StochasticGD
func TestSparseThreeDimensionalLineShouldPass2(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 0, base.L2, 1000, sparseThreeDLineX, threeDLineY, len(threeDLineX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			guess, err = model.Predict([]float64{float64(i), float64(j)})
			assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
			assert.InDelta(t, 10.0+float64(i)/10+float64(j)/5, guess[0], 1e-2, "Guess should be really close to i+x (within 1e-2) for line z=10 + (x+y)/10")
			assert.Nil(t, err, "Prediction error should be nil")
		}
	}
}


// Test Online Learning through channels

func TestSparseOnlineLinearOneDXShouldPass1(t *testing.T) {
	// create the channel of data and errors
	stream := make(chan base.Datapoint, 100)
	errors := make(chan error)

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 0, base.L2, 0, nil, nil, 1)
	go model.OnlineLearn(errors, stream, func(theta [][]float64) {})

	// start passing data to our datastream
	//
	// we could have data already in our channel
	// when we instantiated the Perceptron, though
	for iter := 0; iter < 500; iter++ {
		for i := -40.0; i < 40; i += 0.15 {
			stream <- base.Datapoint{
				X: []float64{i},
				Y: []float64{i/10 + 20},
			}
		}
	}

	// close the dataset
	close(stream)

	err, more := <-errors

	assert.Nil(t, err, "Learning error should be nil")
	assert.False(t, more, "There should be no errors returned")

	// test a larger dataset now
	iter := 0
	for i := -100.0; i < 100; i += 0.347 {
		guess, err := model.Predict([]float64{i})
		assert.Nil(t, err, "Prediction error should be nil")
		assert.Len(t, guess, 1, "Guess should have length 1")

		assert.InDelta(t, i/10+20, guess[0], 1e-2, "Guess should be close to i/10 + 20 for i=%v", i)
		iter++
	}
	fmt.Printf("Iter: %v\n", iter)
}


func TestSparseOnlineLinearOneDXShouldFail1(t *testing.T) {
	// create the channel of data and errors
	stream := make(chan base.Datapoint, 1000)
	errors := make(chan error)

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 0, base.L2, 0, nil, nil, 1)
	go model.OnlineLearn(errors, stream, func(theta [][]float64) {})

	// give invalid data when it should be -1
	for i := -500.0; abs(i) > 1; i *= -0.90 {
		stream <- base.Datapoint{
			X: []float64{i},
			Y: []float64{i/10 + 20, 10, 11},
		}
	}

	// close the dataset
	close(stream)

	count := 0
	for {
		_, more := <-errors
		count++
		if !more {
			assert.True(t, count > 1, "Learning error should not be nil")
			break
		}
	}
}


func TestSparseOnlineLinearOneDXShouldFail2(t *testing.T) {
	// create the channel of data and errors
	stream := make(chan base.Datapoint, 1000)
	errors := make(chan error)

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 0, base.L2, 0, nil, nil, 1)
	go model.OnlineLearn(errors, stream, func(theta [][]float64) {})

	// give invalid data when it should be -1
	for i := -500.0; abs(i) > 1; i *= -0.90 {
		stream <- base.Datapoint{
			X: []float64{i, 0, 13},
			Y: []float64{i/10 + 20},
		}
	}

	// close the dataset
	close(stream)

	count := 0
	for {
		_, more := <-errors
		count++
		if !more {
			assert.True(t, count > 1, "Learning error should not be nil")
			break
		}
	}
}



func TestSparseOnlineLinearOneDXShouldFail3(t *testing.T) {
	// create the channel of errors
	errors := make(chan error)

	model := NewSparseLeastSquares(base.StochasticGD, .0001, .0001, 0, base.L2, 0, nil, nil, 1)
	go model.OnlineLearn(errors, nil, func(theta [][]float64) {})

	err := <-errors
	assert.NotNil(t, err, "Learning error should not be nil")
}


func TestSparseOnlineLinearFourDXShouldPass1(t *testing.T) {
	// create the channel of data and errors
	stream := make(chan base.Datapoint, 100)
	errors := make(chan error)

	var updates int

	model := NewSparseLeastSquares(base.StochasticGD, 1e-5, 1e-5, 0, base.L2, 0, nil, nil, 4)
	go model.OnlineLearn(errors, stream, func(theta [][]float64) {
		updates++
	})

	go func() {
		for iterations := 0; iterations < 25; iterations++ {
			for i := -200.0; abs(i) > 1; i *= -0.75 {
				for j := -200.0; abs(j) > 1; j *= -0.75 {
					for k := -200.0; abs(k) > 1; k *= -0.75 {
						for l := -200.0; abs(l) > 1; l *= -0.75 {
							stream <- base.Datapoint{
								X: []float64{i, j, k, l},
								Y: []float64{i/2 + 2*k - 4*j + 2*l + 3},
							}
						}
					}
				}
			}
		}

		// close the dataset
		close(stream)
	}()

	count := 0
	for {
		err, more := <-errors
		assert.Nil(t, err, "Learning error should be nil")
		count++
		if !more {
			// account (pun intended) for the ++ on every iteration
			//
			// in other words, this should only iterate once, and
			// more should be false in that case
			assert.Equal(t, 0, count-1, "There should be no errors returned")
			break
		}
	}

	assert.True(t, updates > 100, "There should be more than 100 updates of theta")

	for i := -200.0; i < 200; i += 100 {
		for j := -200.0; j < 200; j += 100 {
			for k := -200.0; k < 200; k += 100 {
				for l := -200.0; l < 200; l += 100 {
					guess, err := model.Predict([]float64{i, j, k, l})
					assert.Nil(t, err, "Prediction error should be nil")
					assert.Len(t, guess, 1, "Guess should have length 1")

					assert.InDelta(t, i/2+2*k-4*j+2*l+3, guess[0], 1e-2, "Guess should be close to i/2+2*k-4*j+2*l+3")
				}
			}
		}
	}
}


// Test Persistance To File

// test persisting y=x to file
func TestSparsePersistLeastSquaresShouldPass1(t *testing.T) {
	var err error

	model := NewSparseLeastSquares(base.BatchGD, 1e-6, 1e-6, 0, base.L2, 75, sparseNoisyX, noisyY, len(noisyX[0]))
	err = model.Learn("")
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := 400.0; i < 600; i++ {
		guess, err = model.Predict([]float64{i})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i*0.5, guess[0], 5, "Guess*2 should be close to input for y=0.5*x")
		assert.Nil(t, err, "Prediction error should be nil")
	}

	// not that we know it works, try persisting to file,
	// then resetting the parameter vector theta, then
	// restoring it and testing that predictions are correct
	// again.

	err = model.PersistToFile("/tmp/.goml/LeastSquares.json")
	assert.Nil(t, err, "Persistance error should be nil")

	model.Parameters = make([]float64, len(model.Parameters))

	// make sure it WONT work now that we reset theta
	//
	// the result of Theta transpose * X should always
	// be 0 because theta is the zero vector right now.
	for i := 400.0; i < 600; i++ {
		guess, err = model.Predict([]float64{i})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.Equal(t, 0.0, guess[0], "Guess should be 0 when theta is the zero vector")
		assert.Nil(t, err, "Prediction error should be nil")
	}

	err = model.RestoreFromFile("/tmp/.goml/LeastSquares.json")
	assert.Nil(t, err, "Persistance error should be nil")

	for i := 400.0; i < 600; i++ {
		guess, err = model.Predict([]float64{i})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i*0.5, guess[0], 5, "Guess*2 should be close to input for y=0.5*x")
		assert.Nil(t, err, "Prediction error should be nil")
	}
}
