// This package creates generic functins for crud operations on controller objects
// TODO is this abstraction good? rename package
package adapter

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func GetAll[modelType, newType any](
	name string,
	getData func() ([]modelType, error),
	mapper func(modelType) newType,
) func() ([]newType, error) {
	return func() ([]newType, error) {
		var viewData []newType

		dbData, err := getData()
		if err != nil {
			log.Error(
				fmt.Sprintf("failed to get %s", name),
				"error", err,
			)
			return viewData, err
		}

		for _, dbItem := range dbData {
			viewData = append(viewData, mapper(dbItem))
		}

		return viewData, nil
	}
}

func Get[modelType, newType any](
	name string,
	getData func(int) (modelType, error),
	mapper func(modelType) newType,
) func(int) (newType, error) {
	return func(id int) (newType, error) {
		var viewData newType

		dbData, err := getData(id)
		if err != nil {
			log.Error(
				fmt.Sprintf("failed to get %s", name),
				"error", err,
			)
			return viewData, err
		}

		viewData = mapper(dbData)

		return viewData, nil
	}
}

func Add[modelType, newType any](
	addData func(modelType) (int, error),
	mapper func(newType) modelType,
) func(newType) error {
	return func(viewData newType) error {
		_, err := addData(mapper(viewData))
		return err
	}
}
