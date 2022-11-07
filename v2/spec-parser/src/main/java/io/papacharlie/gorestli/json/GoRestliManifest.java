package io.papacharlie.gorestli.json;

import java.util.HashSet;
import java.util.Set;


public class GoRestliManifest {
  public final String _packageRoot;
  public final Set<DataType> _inputDataTypes = new HashSet<>();
  public final Set<DataType> _dependencyDataTypes = new HashSet<>();
  public final Set<Resource> _resources = new HashSet<>();

  public GoRestliManifest(String packageRoot) {
    _packageRoot = packageRoot;
  }
}
