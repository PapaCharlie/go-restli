package restlicodec

func (e *Encoder) Int32Field(fieldName string, fieldValue int32) {
	e.writeField(fieldName)
	e.encoder.Int32(fieldValue)
}

func (e *Encoder) Int32MapField(fieldName string, fieldValue map[string]int32) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Int32(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) Int32ArrayField(fieldName string, fieldValue []int32) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Int32(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) Int64Field(fieldName string, fieldValue int64) {
	e.writeField(fieldName)
	e.encoder.Int64(fieldValue)
}

func (e *Encoder) Int64MapField(fieldName string, fieldValue map[string]int64) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Int64(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) Int64ArrayField(fieldName string, fieldValue []int64) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Int64(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) Float32Field(fieldName string, fieldValue float32) {
	e.writeField(fieldName)
	e.encoder.Float32(fieldValue)
}

func (e *Encoder) Float32MapField(fieldName string, fieldValue map[string]float32) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Float32(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) Float32ArrayField(fieldName string, fieldValue []float32) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Float32(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) Float64Field(fieldName string, fieldValue float64) {
	e.writeField(fieldName)
	e.encoder.Float64(fieldValue)
}

func (e *Encoder) Float64MapField(fieldName string, fieldValue map[string]float64) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Float64(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) Float64ArrayField(fieldName string, fieldValue []float64) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Float64(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) BoolField(fieldName string, fieldValue bool) {
	e.writeField(fieldName)
	e.encoder.Bool(fieldValue)
}

func (e *Encoder) BoolMapField(fieldName string, fieldValue map[string]bool) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Bool(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) BoolArrayField(fieldName string, fieldValue []bool) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Bool(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) StringField(fieldName string, fieldValue string) {
	e.writeField(fieldName)
	e.encoder.String(fieldValue)
}

func (e *Encoder) StringMapField(fieldName string, fieldValue map[string]string) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.String(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) StringArrayField(fieldName string, fieldValue []string) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.String(item)
	}
	e.encoder.WriteArrayEnd()
}

func (e *Encoder) BytesField(fieldName string, fieldValue []byte) {
	e.writeField(fieldName)
	e.encoder.Bytes(fieldValue)
}

func (e *Encoder) BytesMapField(fieldName string, fieldValue map[string][]byte) {
	e.writeField(fieldName)
	e.encoder.WriteMapStart()
	first := true
	for k, v := range fieldValue {
		if first {
			first = false
		} else {
			e.encoder.WriteMapEntryDelimiter()
		}
		e.encoder.WriteMapKey(k)
		e.encoder.WriteMapKeyDelimiter()
		e.encoder.Bytes(v)
	}
	e.encoder.WriteMapEnd()
}

func (e *Encoder) BytesArrayField(fieldName string, fieldValue [][]byte) {
	e.writeField(fieldName)
	e.encoder.WriteArrayStart()
	for index, item := range fieldValue {
		if index > 0 {
			e.encoder.WriteArrayItemDelimiter()
		}
		e.encoder.Bytes(item)
	}
	e.encoder.WriteArrayEnd()
}
