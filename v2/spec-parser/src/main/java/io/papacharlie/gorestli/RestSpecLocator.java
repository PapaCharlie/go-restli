package io.papacharlie.gorestli;

import com.google.common.collect.ImmutableMap;
import com.linkedin.pegasus.generator.CodeUtil;
import com.linkedin.restli.common.RestConstants;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestSpecCodec;
import java.io.IOException;
import java.net.URI;
import java.nio.file.FileSystem;
import java.nio.file.FileSystems;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;


public class RestSpecLocator {
  private static final RestSpecCodec CODEC = new RestSpecCodec();

  public static Map<Path, ResourceSchema> parseSources(Set<Path> sourcePaths) {
    return sourcePaths.stream()
        .flatMap(RestSpecLocator::walk)
        .flatMap(p -> {
          if (hasExtension(p, ".jar")) {
            return expandJAR(p);
          } else if (hasExtension(p, RestConstants.RESOURCE_MODEL_FILENAME_EXTENSION)) {
            return Stream.of(readSchema(p));
          } else {
            return Stream.of();
          }
        })
        .collect(ImmutableMap.toImmutableMap(p -> p.first, p -> p.second));
  }

  private static CodeUtil.Pair<Path, ResourceSchema> readSchema(Path schema) {
    try {
      return new CodeUtil.Pair<>(schema, CODEC.readResourceSchema(Files.newInputStream(schema)));
    } catch (IOException e) {
      throw new RuntimeException("Failed to read resource schema from " + schema, e);
    }
  }

  private static Stream<CodeUtil.Pair<Path, ResourceSchema>> loadSchemas(Stream<Path> stream) {
    return stream.flatMap(p -> {
      if (hasExtension(p, ".jar")) {
        return expandJAR(p);
      } else if (hasExtension(p, RestConstants.RESOURCE_MODEL_FILENAME_EXTENSION)) {
        return Stream.of(readSchema(p));
      } else {
        return Stream.of();
      }
    });
  }

  private static boolean hasExtension(Path path, String extension) {
    return path.toString().endsWith(extension);
  }

  private static Stream<CodeUtil.Pair<Path, ResourceSchema>> expandJAR(Path jarPath) {
    try (FileSystem fs = FileSystems.newFileSystem(URI.create("jar:" + jarPath.toUri()), Collections.emptyMap())) {
      return loadSchemas(StreamSupport.stream(fs.getRootDirectories().spliterator(), false)
          .flatMap(RestSpecLocator::walk))
          // execute a terminal operation to read all schemas while the jar is open
          .collect(Collectors.toList())
          .stream();
    } catch (IOException e) {
      throw new RuntimeException("Failed to expand " + jarPath, e);
    }
  }

  private static Stream<Path> walk(Path path) {
    try {
      return Files.walk(path);
    } catch (IOException e) {
      throw new RuntimeException("Failed to walk " + path, e);
    }
  }
}
