package io.papacharlie.gorestli;

import com.linkedin.pegasus.generator.CodeUtil.Pair;
import com.linkedin.pegasus.generator.DataSchemaParser;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.tools.clientgen.RestSpecParser;
import com.linkedin.util.FileUtil;
import io.papacharlie.gorestli.json.GoRestliSpec;
import io.papacharlie.gorestli.json.GoRestliSpec.DataType;
import java.io.File;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Set;
import java.util.stream.Collectors;


public class GoRestliRestSpecParser {
  private final String[] _restSpecPaths;
  private final TypeParser _typeParser;
  private final GoRestliSpec _snapshot = new GoRestliSpec();

  public GoRestliRestSpecParser(String resolverPath, String restSpecDir) {
    this(resolverPath, FileUtil.listFiles(new File(restSpecDir), f -> f.getName().endsWith(".restspec.json")).stream()
        .map(File::getAbsolutePath)
        .toArray(String[]::new));
  }

  public GoRestliRestSpecParser(String resolverPath, Set<String> restSpecPaths) {
    this(resolverPath, restSpecPaths.toArray(new String[0]));
  }

  private GoRestliRestSpecParser(String resolverPath, String[] restSpecPaths) {
    _restSpecPaths = restSpecPaths;
    _typeParser = new TypeParser(new DataSchemaParser(resolverPath).getSchemaResolver());
  }

  public GoRestliSpec parse() {
    RestSpecParser.ParseResult restSpecParseResult = new RestSpecParser().parseSources(_restSpecPaths);

    for (Pair<ResourceSchema, File> result : restSpecParseResult.getSchemaAndFiles()) {
      Set<DataType> dataTypes = _typeParser.extractDataTypes(result.first);
      _snapshot._dataTypes.addAll(dataTypes);
      _snapshot._resources.addAll(new ResourceParser(result.first, result.second, _typeParser).parse());
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
