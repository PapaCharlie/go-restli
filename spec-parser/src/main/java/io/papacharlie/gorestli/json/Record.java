package io.papacharlie.gorestli.json;

import com.linkedin.data.ByteString;
import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.json.RestliType.Identifier;
import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.List;

import static io.papacharlie.gorestli.Utils.*;


public class Record extends NamedType {
  public final List<Field> _fields;

  public Record(NamedDataSchema namedDataSchema, Path sourceFile, List<Field> fields) {
    super(namedDataSchema, sourceFile);
    _fields = fields;
  }

  public Field getField(String name) {
    for (Field field : _fields) {
      if (field._name.equals(name)) {
        return field;
      }
    }
    throw new IllegalArgumentException(String.format("No such field %s in %s", name, getIdentifier()));
  }

  public static class Field {
    public final String _name;
    public final String _doc;
    public final RestliType _type;
    public final boolean _isOptional;
    public final String _defaultValue;
    public final Identifier _includedFrom;

    public Field(String name, String doc, RestliType type, Boolean isOptional, Object defaultValue,
        Identifier includedFrom) {
      _name = name;
      _doc = doc;
      _type = type;
      _isOptional = isOptional != null && isOptional;
      _defaultValue = serializeDefaultValue(defaultValue);
      _includedFrom = includedFrom;
    }

    public Field(String name, String doc, RestliType type, Boolean isOptional) {
      this(name, doc, type, isOptional, null, null);
    }

    private static String serializeDefaultValue(Object value) {
      if (value == null) {
        return null;
      }

      if (value instanceof ByteString) {
        return UGLY_GSON.toJson(((ByteString) value).asString(StandardCharsets.UTF_8));
      } else {
        return UGLY_GSON.toJson(value);
      }
    }
  }
}
