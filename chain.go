package chain

import (
	"fmt"
	"reflect"
	"strconv"
)

type Chain struct {
	data interface{}
}

func (c *Chain) Print() *Chain {
	fmt.Println(c.data)
	return c
}

func (c *Chain) CustomFn(f func(param interface{}) interface{}) *Chain {
	c.data = f(c.data)
	return c
}

func (c *Chain) Convert(to string) *Chain {
	switch to {
	case "string":
		c.data = fmt.Sprintf("%v", c.data)
	case "int":
		switch v := c.data.(type) {
		case int:
			// ya es int
		case float64:
			c.data = int(v)
		case string:
			val, err := strconv.Atoi(v)
			if err == nil {
				c.data = val
			} else {
				fmt.Println("No se pudo convertir a int:", err)
			}
		default:
			fmt.Println("Tipo no compatible para convertir a int")
		}
	case "float64":
		switch v := c.data.(type) {
		case float64:
			// ya es float64
		case int:
			c.data = float64(v)
		case string:
			val, err := strconv.ParseFloat(v, 64)
			if err == nil {
				c.data = val
			} else {
				fmt.Println("No se pudo convertir a float64:", err)
			}
		default:
			fmt.Println("Tipo no compatible para convertir a float64")
		}
	case "bool":
		switch v := c.data.(type) {
		case bool:
			// ya es bool
		case string:
			val, err := strconv.ParseBool(v)
			if err == nil {
				c.data = val
			} else {
				fmt.Println("No se pudo convertir a bool:", err)
			}
		default:
			fmt.Println("Tipo no compatible para convertir a bool")
		}
	case "map":
		// Implementa la conversión a map si es necesario
	case "slice":
		// Implementa la conversión a slice si es necesario
	default:
		fmt.Println("Conversión no soportada:", to)
	}
	return c
}

func (c *Chain) SetData(data interface{}) *Chain {
	c.data = data
	return c
}

func (c *Chain) Map(f interface{}) *Chain {
	dataType := reflect.TypeOf(c.data)

	if dataType.Kind() == reflect.Slice {
		slice := reflect.ValueOf(c.data)
		for i := 0; i < slice.Len(); i++ {
			element := slice.Index(i)
			result := reflect.ValueOf(f).Call([]reflect.Value{element})
			slice.Index(i).Set(result[0])
		}
	} else if dataType.Kind() == reflect.Map {
		m := reflect.ValueOf(c.data)
		newMap := reflect.MakeMap(reflect.TypeOf(map[interface{}]interface{}{}))
		for _, key := range m.MapKeys() {
			value := m.MapIndex(key)
			result := reflect.ValueOf(f).Call([]reflect.Value{value})
			newMap.SetMapIndex(key, result[0])
		}
		c.data = newMap.Interface()
	} else {
		fmt.Println("Unsupported data type for mapping:", dataType)
	}

	return c
}

func (c *Chain) Filter(f func(interface{}) bool) *Chain {
	switch v := c.data.(type) {
	case []interface{}:
		var filtered []interface{}
		for _, item := range v {
			if f(item) {
				filtered = append(filtered, item)
			}
		}
		c.data = filtered
	case map[interface{}]interface{}:
		newMap := make(map[interface{}]interface{})
		for k, v := range v {
			if f(v) {
				newMap[k] = v
			}
		}
		c.data = newMap
	default:
		fmt.Println("Tipo de datos no compatible para filtrar:", v)
	}
	return c
}

func New(data interface{}) *Chain {
	return &Chain{data}
}
