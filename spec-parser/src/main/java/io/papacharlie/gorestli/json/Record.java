package io.papacharlie.gorestli.json;

import com.linkedin.data.ByteString;
import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.json.RestliType.Identifier;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.List;

import static io.papacharlie.gorestli.Utils.UGLY_GSON;


public class Record extends NamedType {
  public final List<Identifier> _includes;
  public final List<Field> _fields;

  public Record(NamedDataSchema namedDataSchema, Path sourceFile, List<Identifier> includes, List<Field> fields) {
    super(namedDataSchema, sourceFile);
    _includes = includes;
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
      _isOptional = isOptional != null && isOptional;
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
        return UGLY_GSON.toJson(((ByteString) value).asString(StandardCharsets.UTF_8));
      } else {
        return UGLY_GSON.toJson(value);
      }
    }
  }
}
