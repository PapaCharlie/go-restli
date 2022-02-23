package io.papacharlie.gorestli.json;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.Set;


public class Resource {
  public final String _namespace;
  public final String _doc;
  public final String _sourceFile;
  public final String _rootResourceName;
  public final RestliType _resourceSchema;
  public final Set<String> _readOnlyFields;
  public final Set<String> _createOnlyFields;
  public final List<Method> _methods = new ArrayList<>();
  public final boolean _isCollection;

  public Resource(String namespace, String doc, String sourceFile, String rootResourceName, RestliType resourceSchema,
      Set<String> readOnlyFields, Set<String> createOnlyFields, boolean isCollection) {
    _namespace = namespace;
    _doc = doc;
    _sourceFile = sourceFile;
    _rootResourceName = rootResourceName;
    _resourceSchema = resourceSchema;
    _readOnlyFields = readOnlyFields;
    _createOnlyFields = createOnlyFields;
    _isCollection = isCollection;
  }

  public Resource addMethod(Method m) {
    _methods.add(m);
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
