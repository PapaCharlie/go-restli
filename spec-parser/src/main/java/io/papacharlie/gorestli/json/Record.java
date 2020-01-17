package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.Utils;
import java.io.File;
import java.util.List;


public class Record extends NamedType {
  public final List<Field> _fields;

  public Record(NamedDataSchema namedDataSchema, File sourceFile, List<Field> fields) {
    super(namedDataSchema, sourceFile);
    _fields = fields;
  }

  public static class Field {
    public final String _name;
    public final String _doc;
    public final RestliType _type;
    public final boolean _isOptional;
    public final String _defaultValue;

    public Field(String name, String doc, RestliType type, Boolean isOptional, Object defaultValue) {
      _name = name;
      _doc = doc;
      _type = type;
      _isOptional = (isOptional == null) ? false : isOptional;
      _defaultValue = (defaultValue == null) ? null : Utils.toJson(defaultValue);
    }

    public Field(String name, String doc, RestliType type, Boolean isOptional) {
      this(name, doc, type, isOptional, null);
    }
  }
}
