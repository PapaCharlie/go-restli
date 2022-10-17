package io.papacharlie.gorestli.json;

import java.util.HashSet;
import java.util.Set;


public class GoRestliManifest {
  public final String _packageRoot;
  public final Set<DataType> _dataTypes = new HashSet<>();
  public final Set<DataType> _additionalDataTypes = new HashSet<>();
  public final Set<Resource> _resources = new HashSet<>();

  public GoRestliManifest(String packageRoot) {
    _packageRoot = packageRoot;
  }
}
