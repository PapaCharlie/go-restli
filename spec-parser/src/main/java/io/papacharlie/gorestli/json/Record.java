package io.papacharlie.gorestli.json;

import com.linkedin.data.ByteString;
import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.Utils;
import java.io.File;
import java.nio.charset.StandardCharsets;
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
      _defaultValue = serializeDefaultValue(defaultValue);
    }

    public Field(String name, String doc, RestliType type, Boolean isOptional) {
      this(name, doc, type, isOptional, null);
    }

    private static String serializeDefaultValue(Object value) {
      if (value == null) {
        return null;
      }

      if (value instanceof ByteString) {
        return Utils.toJson(((ByteString) value).asString(StandardCharsets.UTF_8));
      } else {
        return Utils.toJson(value);
      }
    }
  }
}
