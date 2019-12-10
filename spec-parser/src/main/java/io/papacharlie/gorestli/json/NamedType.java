package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import java.io.File;
import java.util.Objects;


public abstract class NamedType {
  public final String _name;
  public final String _namespace;
  public final String _doc;
  public final String _sourceFile;

  protected NamedType(NamedDataSchema namedDataSchema, File sourceFile) {
    this(namedDataSchema.getName(), namedDataSchema.getNamespace(), namedDataSchema.getDoc(), sourceFile);
  }

  protected NamedType(String name, String namespace, String doc, File sourceFile) {
    _name = name;
    _namespace = namespace;
    _doc = doc;
    _sourceFile = sourceFile.getAbsolutePath();
  }

  @Override
  public int hashCode() {
    return Objects.hash(_namespace, _name);
  }

  @Override
  public boolean equals(Object obj) {
    if (!(obj instanceof NamedType)) {
      return false;
    }
    NamedType namedType = (NamedType) obj;
    return _namespace.equals(namedType._namespace) && _name.equals(namedType._name);
  }
}
