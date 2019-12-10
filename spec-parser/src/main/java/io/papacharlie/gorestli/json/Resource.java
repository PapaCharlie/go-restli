package io.papacharlie.gorestli.json;

import io.papacharlie.gorestli.Utils;
import java.util.List;
import java.util.Objects;


public class Resource {
  public final String _namespace;
  public final String _doc;
  public final String _sourceFile;
  public final String _rootResourceName;
  public final RestliType _resourceSchema;
  public List<Method> _methods;

  public Resource(String namespace, String doc, String sourceFile, String rootResourceName, RestliType resourceSchema) {
    _namespace = namespace;
    _doc = doc;
    _sourceFile = sourceFile;
    _rootResourceName = rootResourceName;
    _resourceSchema = resourceSchema;
  }

  public Resource addMethod(Method m) {
    _methods = Utils.append(_methods, m);
    return this;
  }


  @Override
  public int hashCode() {
    return Objects.hash(_namespace);
  }

  @Override
  public boolean equals(Object obj) {
    if (!(obj instanceof Resource)) {
      return false;
    }
    Resource other = (Resource) obj;
    return _namespace.equals(other._namespace);
  }
}
