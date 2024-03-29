package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.TyperefDataSchema;
import io.papacharlie.gorestli.json.RestliType.GoPrimitive;
import java.nio.file.Path;
import java.util.Map;


public class Typeref extends NamedType {
  public final GoPrimitive _type;
  public final Map<String, Object> _properties;

  public Typeref(TyperefDataSchema schema, Path sourceFile, GoPrimitive type) {
    super(schema, sourceFile);
    _type = type;
    _properties = schema.getMergedTyperefProperties();
  }
}
