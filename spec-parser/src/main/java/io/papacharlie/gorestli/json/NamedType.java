package io.papacharlie.gorestli.json;

import com.linkedin.data.schema.NamedDataSchema;
import io.papacharlie.gorestli.json.RestliType.Identifier;
import java.io.File;


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

  public Identifier getIdentifier() {
    return new Identifier(_namespace, _name);
  }

  public RestliType restliType() {
    return new RestliType(getIdentifier());
  }

  @Override
  public int hashCode() {
    return getIdentifier().hashCode();
  }

  @Override
  public boolean equals(Object obj) {
    if (!(obj instanceof NamedType)) {
      return false;
    }
    NamedType namedType = (NamedType) obj;
    return getIdentifier().equals(namedType.getIdentifier());
  }
}
