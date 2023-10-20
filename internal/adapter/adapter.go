// This package creates generic functins for crud operations on controller objects
// TODO is this abstraction good? rename package
package adapter

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func GetAll[modelType, viewType any](
	name string,
	getData func() ([]modelType, error),
	mapper func(modelType) viewType,
) func() ([]viewType, error) {
	return func() ([]viewType, error) {
		var viewData []viewType

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

func Get[modelType, viewType any](
	name string,
	getData func(int) (modelType, error),
	mapper func(modelType) viewType,
) func(int) (viewType, error) {
	return func(id int) (viewType, error) {
		var viewData viewType

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

func Add[modelType, viewType any](
	addData func(modelType) (int, error),
	mapper func(viewType) modelType,
) func(viewType) (int, error) {
	return func(viewData viewType) (int, error) {
		return addData(mapper(viewData))
	}
}
