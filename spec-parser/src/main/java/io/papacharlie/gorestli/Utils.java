package io.papacharlie.gorestli;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonDeserializer;
import com.google.gson.JsonPrimitive;
import com.google.gson.JsonSerializer;
import com.linkedin.restli.restspec.CustomAnnotationContentSchema;
import com.linkedin.restli.restspec.ResourceSchema;
import com.linkedin.restli.restspec.RestMethodSchema;
import io.papacharlie.gorestli.json.RestliType.GoPrimitive;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeFormatterBuilder;
import java.time.temporal.ChronoField;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import org.apache.commons.lang3.StringUtils;


public class Utils {
  public static final Gson UGLY_GSON = new GsonBuilder()
      .setFieldNamingStrategy(f -> StringUtils.removeStart(f.getName(), "_"))
      .registerTypeAdapter(GoPrimitive.class,
          (JsonSerializer<GoPrimitive>) (src, typeOfSrc, context) -> new JsonPrimitive(src.name().toLowerCase()))
      .registerTypeAdapter(Path.class,
          (JsonSerializer<Path>) (src, typeOfSrc, context) -> new JsonPrimitive(src.toString()))
      .registerTypeAdapter(Path.class,
          (JsonDeserializer<Path>) (src, typeOfSrc, context) -> Paths.get(src.getAsString()))
      .create();
  public static final Gson GSON = UGLY_GSON.newBuilder().setPrettyPrinting().create();
  private static final DateTimeFormatter LOG_TIME_FORMAT = new DateTimeFormatterBuilder()
      .appendValue(ChronoField.YEAR, 4)
      .appendLiteral('/')
      .appendValue(ChronoField.MONTH_OF_YEAR, 2)
      .appendLiteral('/')
      .appendValue(ChronoField.DAY_OF_MONTH, 2)
      .appendLiteral(' ')
      .appendValue(ChronoField.HOUR_OF_DAY, 2)
      .appendLiteral(':')
      .appendValue(ChronoField.MINUTE_OF_HOUR, 2)
      .appendLiteral(':')
      .appendValue(ChronoField.SECOND_OF_MINUTE, 2)
      .appendLiteral('.')
      .appendValue(ChronoField.MILLI_OF_SECOND, 3)
      .toFormatter();

  private Utils() { /* No instance for you */ }

  public static <T> String toJson(T obj) {
    return GSON.toJson(obj);
  }

  public static <T> List<T> append(List<T> original, T newValue) {
    List<T> newList = new ArrayList<>(emptyIfNull(original));
    newList.add(newValue);
    return newList;
  }

  public static <T> Set<T> emptyIfNull(Set<T> set) {
    return (set == null)
        ? Collections.emptySet()
        : set;
  }

  public static <T> List<T> emptyIfNull(List<T> list) {
    return (list == null)
        ? Collections.emptyList()
        : list;
  }

  public static void log(String format, Object... args) {
    System.err.printf(LOG_TIME_FORMAT.format(LocalDateTime.now()) + " [go-restli] " + format + "%n", args);
  }

  private static final String FORCED_EXPORT_PREFIX = "Exported_";

  public static String exportedIdentifier(String identifier) {
    StringBuilder buf = new StringBuilder();

    int firstCodePoint = identifier.codePointAt(0);
    if (!Character.isAlphabetic(firstCodePoint)) {
      buf.append(FORCED_EXPORT_PREFIX);
      if (identifier.charAt(0) == '_') {
        identifier = identifier.substring(1);
      }
    } else {
      buf.appendCodePoint(Character.toUpperCase(firstCodePoint));
      identifier = identifier.substring(1);
    }

    for (int i = 0; i < identifier.length(); i++) {
      int c = identifier.codePointAt(i);
      if (Character.isAlphabetic(c) || Character.isDigit(c)) {
        buf.appendCodePoint(c);
      } else {
        buf.append('_');
      }
    }

    return buf.toString();
  }

  public static boolean supportsReturnEntity(RestMethodSchema schema) {
    if (schema.getAnnotations() == null) {
      return false;
    }
    return schema.getAnnotations().get("returnEntity") != null;
  }

  @SuppressWarnings("unchecked")
  public static Set<String> extractAnnotationValues(ResourceSchema schema, String key) {
    if (schema.getAnnotations() == null) {
      return null;
    }

    CustomAnnotationContentSchema customAnnotations = schema.getAnnotations().get(key);
    if (customAnnotations == null) {
      return null;
    }

    return new HashSet<>((List<String>) customAnnotations.data().get("value"));
  }
}
