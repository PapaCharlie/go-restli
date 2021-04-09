package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.util.Objects;


public class Identifier {
  public final String _namespace;
  public final String _name;

  public Identifier(String namespace, String name) {
    _namespace = namespace;
    _name = name;
  }

  public Identifier(NamedDataSchema schema) {
    this(schema.getNamespace(), schema.getName());
  }

  @Override
  public int hashCode() {
    return Objects.hash(_namespace, _name);
  }

  @Override
  public boolean equals(Object obj) {
    if (!(obj instanceof Identifier)) {
      return false;
    }
    Identifier identifier = (Identifier) obj;
    return _namespace.equals(identifier._namespace)
        && _name.equals(identifier._name);
  }

  @Override
  public String toString() {
    return _namespace + "." + _name;
  }
}
