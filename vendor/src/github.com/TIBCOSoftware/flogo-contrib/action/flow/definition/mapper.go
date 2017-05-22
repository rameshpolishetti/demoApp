package definition

import (
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)


// MapperDef represents a Mapper, which is a collection of mappings
type MapperDef struct {
	//todo possibly add optional lang/mapper type so we can fast fail on unsupported mappings/mapper combo
	Mappings []*data.MappingDef
}

type MapperFactory interface {

	// NewMapper creates a new Mapper from the specified MapperDef
	NewMapper(mapperDef *MapperDef) data.Mapper

	// NewTaskInputMapper creates a new Input Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// NewTaskOutputMapper creates a new Output Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// GetDefaultTaskOutputMapper get the default Output Mapper for the
	// specified Task
	GetDefaultTaskOutputMapper(task *Task) data.Mapper
}

var	mapperFactory MapperFactory

func SetMapperFactory(mapper MapperFactory) {
	mapperFactory = mapper
}

func GetMapperFactory() MapperFactory {

	//temp hack until we fully move to new flow action model
	if mapperFactory == nil {
		mapperFactory = &BasicMapperFactory{}
	}

	return mapperFactory
}


//todo move the following to flowAction

type BasicMapperFactory struct {

}

func(mf *BasicMapperFactory) NewMapper(mapperDef *MapperDef) data.Mapper {
	return NewBasicMapper(mapperDef)
}

func(mf *BasicMapperFactory) NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	return NewBasicMapper(mapperDef)
}

func(mf *BasicMapperFactory) NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	return NewBasicMapper(mapperDef)
}

func (mf *BasicMapperFactory) GetDefaultTaskOutputMapper(task *Task) data.Mapper {
	return &DefaultOutputMapper{task:task}
}

// BasicMapper is a simple object holding and executing mappings
type BasicMapper struct {
	mappings []*data.MappingDef
}

// NewBasicMapper creates a new BasicMapper with the specified mappings
func NewBasicMapper(mapperDef *MapperDef) data.Mapper {

	var mapper BasicMapper
	mapper.mappings = mapperDef.Mappings

	return &mapper
}

// Mappings gets the mappings of the BasicMapper
func (m *BasicMapper) Mappings() []*data.MappingDef {
	return m.mappings
}

// Apply executes the mappings using the values from the input scope
// and puts the results in the output scope
//
// return error
func (m *BasicMapper) Apply(inputScope data.Scope, outputScope data.Scope) {

	//todo validate types
	for _, mapping := range m.mappings {

		switch mapping.Type {
		case data.MtAssign:

			attrName, attrPath, pathType := data.GetAttrPath(mapping.Value)

			tv, exists := inputScope.GetAttr(attrName)
			attrValue := tv.Value

			if exists && len(attrPath) > 0 {
				if tv.Type == data.PARAMS {
					valMap := attrValue.(map[string]string)
					attrValue, exists = valMap[attrPath]
				} else if tv.Type == data.ARRAY && pathType == data.PT_ARRAY {
					//assigning part of array
					idx, _ := strconv.Atoi(attrPath)
					//todo handle err
					valArray := attrValue.([]interface{})
					attrValue = valArray[idx]
				} else {
					//for now assume if we have a path, attr is "object"
					valMap := attrValue.(map[string]interface{})
					attrValue = data.GetMapValue(valMap, attrPath)
					//attrValue, exists = valMap[attrPath]
				}
			}

			//todo implement type conversion
			if exists {

				attrName, attrPath, pathType := data.GetAttrPath(mapping.MapTo)
				toAttr, oe := outputScope.GetAttr(attrName)

				if !oe {
					//todo handle attr dne
					fmt.Printf("Attr %s not found in output scope\n", attrName)
					return
				}

				switch pathType {
				case data.PT_SIMPLE:
					outputScope.SetAttrValue(mapping.MapTo, attrValue)
				case data.PT_ARRAY:
					if toAttr.Type == data.ARRAY {
						var valArray []interface{}
						if toAttr.Value == nil {
							//what should we do in this case, construct the array?
							//valArray = make(map[string]string)
						} else {
							valArray = toAttr.Value.([]interface{})
						}

						idx, _ := strconv.Atoi(attrPath)
						//todo handle err
						valArray[idx] = attrValue

						outputScope.SetAttrValue(attrName, valArray)
					} else {
						//todo throw error.. not an ARRAY
					}
				case data.PT_MAP:

					if toAttr.Type == data.PARAMS {
						var valMap map[string]string
						if toAttr.Value == nil {
							valMap = make(map[string]string)
						} else {
							valMap = toAttr.Value.(map[string]string)
						}
						strVal, _ := data.CoerceToString(attrValue)
						valMap[attrPath] = strVal

						outputScope.SetAttrValue(attrName, valMap)
					} else if toAttr.Type == data.OBJECT {
						var valMap map[string]interface{}
						if toAttr.Value == nil {
							valMap = make(map[string]interface{})
						} else {
							valMap = toAttr.Value.(map[string]interface{})
						}
						valMap[attrPath] = attrValue

						outputScope.SetAttrValue(attrName, valMap)
					} else {
						//todo throw error.. not a OBJECT or PARAMS
					}
				}
			}
		//todo: should we ignore if DNE - if we have to add dynamically what type do we use
		case data.MtLiteral:
			outputScope.SetAttrValue(mapping.MapTo, mapping.Value)
		case data.MtExpression:
		//todo implement script mapping
		}
	}
}

// BasicMapper is a simple object holding and executing mappings
type DefaultOutputMapper struct {
	task *Task
}

func (m *DefaultOutputMapper) Apply(inputScope data.Scope, outputScope data.Scope) {

	oscope := outputScope.(data.MutableScope)

	act := activity.Get(m.task.ActivityRef())

	attrNS := "{A" + strconv.Itoa(m.task.ID()) + "."

	for _, attr := range act.Metadata().Outputs {

		oAttr, _ := inputScope.GetAttr(attr.Name)

		if oAttr != nil {
			oscope.AddAttr(attrNS+attr.Name+"}", attr.Type, oAttr.Value)
		}
	}
}
