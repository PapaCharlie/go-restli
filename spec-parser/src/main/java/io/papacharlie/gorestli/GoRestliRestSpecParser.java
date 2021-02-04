package io.papacharlie.gorestli;

import com.google.common.base.Preconditions;
import com.linkedin.data.schema.NamedDataSchema;
import com.linkedin.pegasus.generator.CodeUtil.Pair;
import com.linkedin.pegasus.generator.DataSchemaParser;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.tools.clientgen.RestSpecParser;
import com.linkedin.restli.tools.snapshot.gen.SnapshotGenerator;
import io.papacharlie.gorestli.json.GoRestliSpec;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.Reader;
import java.io.UncheckedIOException;
import java.util.Collections;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;


public class GoRestliRestSpecParser {
  private GoRestliRestSpecParser() { }

  /**
   * Extract a {@link GoRestliSpec} from the given rest specs and PDSC/PDL files
   * @param resolverPath The root directory that contains all the PDSC/PDL schema definitions
   * @param restSpecPaths The set of rest specs to generate bindings for. Cannot be empty unless
   *                      {@code namedDataSchemasToGenerate} is not empty.
   * @param namedDataSchemasToGenerate A set of named types to generate bindings for.
   * @param rawRecords Types that appear in this set will be treated as raw record types and be presented as a
   *                   {@code map[string]interface{}}
   */
  public static GoRestliSpec parse(String resolverPath, Set<String> restSpecPaths,
      Set<String> namedDataSchemasToGenerate, Set<String> rawRecords) {
    if (restSpecPaths.isEmpty() && namedDataSchemasToGenerate.isEmpty()) {
      throw new IllegalArgumentException("Must specify at least one rest spec or a named data schema to generate!");
    }

    GoRestliSpec parsedSpec = new GoRestliSpec();
    DataSchemaParser dataSchemaParser = new DataSchemaParser(resolverPath);
    Map<String, NamedDataSchema> schemas = loadAllSchemas(dataSchemaParser);
    TypeParser typeParser = new TypeParser(dataSchemaParser, rawRecords);

    RestSpecParser.ParseResult restSpecParseResult =
        new RestSpecParser().parseSources(restSpecPaths.toArray(new String[0]));

    for (Pair<ResourceSchema, File> result : restSpecParseResult.getSchemaAndFiles()) {
      typeParser.extractDataTypes(result.first);
      parsedSpec._resources.addAll(new ResourceParser(result.first, result.second, typeParser).parse());
    }

    Set<NamedDataSchema> extraSchemas = new HashSet<>();
    for (String schemaName : namedDataSchemasToGenerate) {
      NamedDataSchema schema = schemas.get(schemaName);
      Preconditions.checkState(schema != null, "%s was not found in %s", schemaName, resolverPath);
      extraSchemas.addAll(expandSchema(dataSchemaParser, schema));
    }

    for (NamedDataSchema schema : extraSchemas) {
      typeParser.addNamedDataSchema(schema);
    }
    parsedSpec._dataTypes.addAll(typeParser.getDataTypes());

    return parsedSpec;
  }

  /**
   * Returns the set of all {@link NamedDataSchema}s that the given schema depends on (including the given schema
   * itself)
   */
  private static Set<NamedDataSchema> expandSchema(DataSchemaParser parser, NamedDataSchema schema) {
    ResourceSchema resourceSchema = new ResourceSchema()
        .setSchema(schema.getFullName());
    // Author's note: creating a fake ResourceSchema just to use the SnapshotGenerator may seem like hacky but it's the
    // best way to ensure logical parity
    SnapshotGenerator snapshotGenerator = new SnapshotGenerator(resourceSchema, parser.getSchemaResolver());
    return new HashSet<>(snapshotGenerator.generateModelList());
  }

  /**
   * Read all the schemas provided in the {@link DataSchemaParser#getResolverPath()}
   */
  private static Map<String, NamedDataSchema> loadAllSchemas(DataSchemaParser parser) {
    DataSchemaParser.ParseResult res;
    try {
      res = parser.parseSources(new String[]{parser.getResolverPath()});
    } catch (IOException e) {
      throw new UncheckedIOException(e);
    }
    return res.getSchemaAndLocations().keySet().stream()
        .filter(s -> s instanceof NamedDataSchema)
        .map(s -> (NamedDataSchema) s)
        .collect(Collectors.toMap(NamedDataSchema::getFullName, Function.identity()));
  }

  private static class StdinParameters {
    String _resolverPath;
    Set<String> _restSpecPaths;
    Set<String> _namedDataSchemasToGenerate;
    Set<String> _rawRecords;
  }

  public static void main(String[] args) throws IOException {
    if (System.in == null) {
      throw new IllegalStateException("Must supply parameters via stdin JSON");
    }

    StdinParameters parameters;
    try (Reader reader = new InputStreamReader(System.in)) {
      parameters = Utils.GSON.fromJson(reader, StdinParameters.class);
    }
    parameters._restSpecPaths = (parameters._restSpecPaths == null)
        ? Collections.emptySet()
        : parameters._restSpecPaths;
    parameters._namedDataSchemasToGenerate = (parameters._namedDataSchemasToGenerate == null)
        ? Collections.emptySet()
        : parameters._namedDataSchemasToGenerate;
    parameters._rawRecords = (parameters._rawRecords == null)
        ? Collections.emptySet()
        : parameters._rawRecords;

    GoRestliSpec spec = parse(
        parameters._resolverPath,
        parameters._restSpecPaths,
        parameters._namedDataSchemasToGenerate,
        parameters._rawRecords);
    System.out.println(Utils.toJson(spec));
  }
}
