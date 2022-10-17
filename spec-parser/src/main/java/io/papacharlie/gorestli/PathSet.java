package io.papacharlie.gorestli;

import java.nio.file.Path;
import java.util.Collection;
import java.util.Map;
import java.util.TreeMap;


public class PathSet {
  private final Map<String, PathSet> _segments = new TreeMap<>();
  private boolean _leaf;

  public PathSet() {
  }

  public PathSet(Collection<Path> paths) {
    if (paths != null) {
      paths.forEach(this::add);
    }
  }

  public void add(Path path) {
    PathSet set = this;
    for (Path segment : path.toAbsolutePath()) {
      set = set._segments.computeIfAbsent(segment.toString(), k -> new PathSet());
    }
    set._leaf = true;
  }

  public boolean contains(Path path) {
    PathSet set = this;
    for (Path segment : path.toAbsolutePath()) {
      PathSet child = set._segments.get(segment.toString());
      if (child != null) {
        set = child;
      } else {
        break;
      }
    }

    return set._leaf;
  }
}
