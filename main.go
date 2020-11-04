package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msgf("usage: %s <file> [<file>...]", os.Args[0])
	}
	var res interface{}
	for _, f := range os.Args[1:] {
		bs, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read file")
		}
		var part interface{}
		err = yaml.Unmarshal(bs, &part)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse file")
		}
		res, err = Merge(res, part)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to merge file")
		}
	}
	bs, err := yaml.Marshal(res)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal result")
	}
	fmt.Println(string(bs))
}

func Merge(a, b interface{}) (_ interface{}, err error) {
	log.Info().Msgf("merge %v (%T) %v (%T)", a, a, b, b)
	switch typedA := a.(type) {
	case []interface{}:
		typedB, ok := b.([]interface{})
		if !ok {
			return nil, errors.New("wrong type on right side")
		}
		return append(typedA, typedB...), nil
	case map[interface{}]interface{}:
		typedB, ok := b.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("wrong type on right side")
		}
		for key, rightVal := range typedB {
			leftVal, ok := typedA[key]
			if !ok {
				typedA[key] = rightVal
			} else {
				typedA[key], err = Merge(leftVal, rightVal)
				if err != nil {
					return nil, err
				}
			}
		}
		return typedA, nil
	default:
		return b, nil
	}
	return nil, errors.New("unexpected end")
}
