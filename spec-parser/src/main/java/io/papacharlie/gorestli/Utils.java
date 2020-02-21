package io.papacharlie.gorestli;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeFormatterBuilder;
import java.time.temporal.ChronoField;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import org.apache.commons.lang3.StringUtils;


public class Utils {
  private static final Gson GSON = new GsonBuilder()
      .setFieldNamingStrategy(f -> StringUtils.removeStart(f.getName(), "_"))
      .setPrettyPrinting()
      .create();
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

  public static <T> List<T> emptyIfNull(List<T> list) {
    return (list == null)
        ? Collections.emptyList()
        : list;
  }

  public static void log(String format, Object... args) {
    System.err.printf("[go-restli] " + LOG_TIME_FORMAT.format(LocalDateTime.now()) + " " + format, args);
  }
}
