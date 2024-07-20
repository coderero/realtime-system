package main

import "github.com/coderero/toll_calculator/types"

type CalculatorService interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type calculatorService struct {
}

func NewCalculatorService() CalculatorService {
	return &calculatorService{}
}

func (c *calculatorService) CalculateDistance(obuData types.OBUData) (float64, error) {
	return 0, nil
}
