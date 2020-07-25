package io.papacharlie.gorestli;

import com.linkedin.pegasus.generator.CodeUtil.Pair;
import com.linkedin.pegasus.generator.DataSchemaParser;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.tools.clientgen.RestSpecParser;
import com.linkedin.util.FileUtil;
import io.papacharlie.gorestli.json.GoRestliSpec;
import java.io.File;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Set;
import java.util.stream.Collectors;


public class GoRestliRestSpecParser {
  private final String _resolverPath;
  private final Set<String> _restSpecPaths;
  private final GoRestliSpec _snapshot = new GoRestliSpec();

  public GoRestliRestSpecParser(String resolverPath, String restSpecDir) {
    this(resolverPath, FileUtil.listFiles(new File(restSpecDir), f -> f.getName().endsWith(".restspec.json")).stream()
        .map(File::getAbsolutePath)
        .collect(Collectors.toSet()));
  }

  public GoRestliRestSpecParser(String resolverPath, Set<String> restSpecPaths) {
    _resolverPath = resolverPath;
    _restSpecPaths = restSpecPaths;
  }

  public GoRestliSpec parse() {
    RestSpecParser.ParseResult restSpecParseResult =
        new RestSpecParser().parseSources(_restSpecPaths.toArray(new String[0]));

    for (Pair<ResourceSchema, File> result : restSpecParseResult.getSchemaAndFiles()) {
      TypeParser typeParser = new TypeParser(new DataSchemaParser(_resolverPath).getSchemaResolver());
      typeParser.extractDataTypes(result.first);
      _snapshot._resources.addAll(new ResourceParser(result.first, result.second, typeParser).parse());
      _snapshot._dataTypes.addAll(typeParser.getDataTypes());
    }

    return _snapshot;
  }

  public static void main(String[] args) {
    if (args.length < 2) {
      System.err.println("Usage: PEGASUS_DIR REST_SPEC,[REST_SPEC...]");
      System.exit(1);
    }

    String resolverPath = Paths.get(args[0]).toAbsolutePath().normalize().toString();
    Set<String> restSpecs = Arrays.stream(Arrays.copyOfRange(args, 1, args.length))
        .map(p -> Paths.get(p).toAbsolutePath().normalize().toString())
        .peek(System.err::println)
        .collect(Collectors.toSet());
    System.out.println(Utils.toJson(new GoRestliRestSpecParser(resolverPath, restSpecs).parse()));
  }
}
