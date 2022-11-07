package io.papacharlie.gorestli.json;

import javax.annotation.Nullable;


public class ResourcePathSegment {
  public final String _resourceName;
  @Nullable
  public final PathKey _pathKey;

  public ResourcePathSegment(String resourceName, @Nullable PathKey pathKey) {
    _resourceName = resourceName;
    _pathKey = pathKey;
  }
}
