package io.papacharlie.gorestli;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.Sets;
import com.linkedin.data.schema.NamedDataSchema;
import com.linkedin.pegasus.generator.DataSchemaParser;
import com.linkedin.restli.restspec.ResourceSchema;
import io.papacharlie.gorestli.json.GoRestliManifest;
import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;
import org.apache.commons.io.IOUtils;


public class GoRestliRestSpecParser {
  private GoRestliRestSpecParser() { }

  /**
   * Extract a {@link GoRestliManifest} from the given rest specs and PDSC/PDL files
   * @param packageRoot The target Golang package that the files should be generated under.
   * @param dependencies The directories, files or JARs that contains all the PDSC/PDL schema definitions required to
   *                     generate the inputs.
   * @param inputs The directories, files or JARs that contain the PDSC/PDL and restspec.json files to generate code
   *                for. Cannot be empty.
   * @param rawRecords Types that appear in this set will be treated as raw record types and be presented as a
   *                   {@code map[string]interface{}}
   * @return The parsed manifest
   * @throws IOException When schemas cannot be parsed
   */
  public static GoRestliManifest parse(
      String packageRoot,
      Set<Path> dependencies,
      Set<Path> inputs,
      Set<String> rawRecords
  ) throws IOException {
    if (inputs.isEmpty()) {
      throw new IllegalArgumentException("Must provide at least one input!");
    }
    Utils.log("Loading schema definitions");

    GoRestliManifest manifest = new GoRestliManifest(packageRoot);
    DataSchemaParser dataSchemaParser = new DataSchemaParser(Sets.union(dependencies, inputs).stream()
        .map(p -> p.toAbsolutePath().toString())
        .collect(Collectors.joining(File.pathSeparator)));

    // Initialize the DataSchemaParser with the dependencies and inputs
    Map<String, NamedDataSchema> allSchemas = loadAllSchemas(dataSchemaParser, Sets.union(dependencies, inputs));

    Map<Path, ResourceSchema> restSpecs = RestSpecLocator.parseSources(inputs);

    if (allSchemas.isEmpty() && restSpecs.isEmpty()) {
      throw new IllegalArgumentException("Must specify at least one rest spec or a named data schema to generate!");
    }

    TypeParser typeParser = new TypeParser(dataSchemaParser, rawRecords);
    typeParser.addNamedDataSchemas(allSchemas);

    Utils.log("Parsed all data models");

    for (Map.Entry<Path, ResourceSchema> entry : restSpecs.entrySet()) {
      manifest._resources.addAll(ResourceParser.parse(entry.getValue(), entry.getKey(), typeParser));
    }

    Utils.log("Parsed all definitions");

    Set<GoRestliManifest.DataType> inputTypes = typeParser.getDataTypes();
    PathSet inputPaths = new PathSet(inputs);
    inputTypes.removeIf(dt -> !inputPaths.contains(dt.getNamedType()._sourceFile));
    manifest._dataTypes.addAll(inputTypes);

    Utils.log("Completed manifest");

    return manifest;
  }

  /**
   * Read all the schemas provided in {@code sources}
   */
  private static Map<String, NamedDataSchema> loadAllSchemas(DataSchemaParser parser, Set<Path> sources)
      throws IOException {
    if (sources.isEmpty()) {
      return Collections.emptyMap();
    }

    Set<Path> absoluteSources = sources.stream()
        .map(Path::toAbsolutePath)
        .collect(Collectors.toSet());

    DataSchemaParser.ParseResult res = parser.parseSources(absoluteSources.stream()
        .map(Path::toString)
        .toArray(String[]::new));

    return res.getSchemaAndLocations().keySet().stream()
        .filter(s -> s instanceof NamedDataSchema)
        .map(s -> (NamedDataSchema) s)
        .collect(ImmutableMap.toImmutableMap(NamedDataSchema::getFullName, Function.identity()));
  }

  private static class StdinParameters {
    String _packageRoot;
    Set<Path> _dependencies;
    Set<Path> _inputs;
    Set<String> _rawRecords;
  }

  public static void main(String[] args) throws IOException {
    String input = IOUtils.toString(System.in);
    System.err.println(input);

    StdinParameters parameters = Utils.GSON.fromJson(input, StdinParameters.class);
    GoRestliManifest manifest = parse(
        parameters._packageRoot,
        Utils.emptyIfNull(parameters._dependencies),
        Utils.emptyIfNull(parameters._inputs),
        Utils.emptyIfNull(parameters._rawRecords)
    );
    System.out.println(Utils.toJson(manifest));
  }
}
