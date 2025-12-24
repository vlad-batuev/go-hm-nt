package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(v any) error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Проверяем, что передан указатель или структура
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("Validate ожидает структуру, получен %s", val.Kind())
	}

	// Проходим по всем полям структуры
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		tag := field.Tag.Get("validate")

		// Пропускаем поля без тега validate
		if tag == "" {
			continue
		}

		// Разбиваем тег на отдельные правила
		rules := strings.Split(tag, ";")

		// Проверяем каждое правило для поля
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}

			// Разбиваем правило на ключ и значение
			parts := strings.SplitN(rule, "=", 2)
			if len(parts) != 2 {
				continue // Пропускаем некорректные правила
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Проверяем правило в зависимости от типа поля
			switch fieldValue.Kind() {
			case reflect.String:
				strValue := fieldValue.String()
				if err := validateStringField(strValue, key, value, field.Name); err != nil {
					return err
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue := fieldValue.Int()
				if err := validateIntField(intValue, key, value, field.Name); err != nil {
					return err
				}

			default:
				// Для других типов пропускаем проверку
				continue
			}
		}
	}

	return nil
}

// validateStringField проверяет строковые поля
func validateStringField(value, ruleKey, ruleValue, fieldName string) error {
	switch ruleKey {
	case "min":
		min, err := strconv.Atoi(ruleValue)
		if err != nil {
			return fmt.Errorf("неверный формат min для поля %s: %v", fieldName, err)
		}
		// Используем руны для корректного подсчета символов (включая юникод)
		if len([]rune(value)) < min {
			return fmt.Errorf("поле %s должно содержать минимум %d символов", fieldName, min)
		}

	case "max":
		max, err := strconv.Atoi(ruleValue)
		if err != nil {
			return fmt.Errorf("неверный формат max для поля %s: %v", fieldName, err)
		}
		if len([]rune(value)) > max {
			return fmt.Errorf("поле %s должно содержать максимум %d символов", fieldName, max)
		}

	case "regexp":
		// Проверяем регулярное выражение
		matched, err := regexp.MatchString(ruleValue, value)
		if err != nil {
			return fmt.Errorf("некорректное регулярное выражение для поля %s: %v", fieldName, err)
		}
		if !matched {
			return fmt.Errorf("поле %s не соответствует формату", fieldName)
		}
	}

	return nil
}

// validateIntField проверяет целочисленные поля
func validateIntField(value int64, ruleKey, ruleValue, fieldName string) error {
	switch ruleKey {
	case "min":
		min, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return fmt.Errorf("неверный формат min для поля %s: %v", fieldName, err)
		}
		if value < min {
			return fmt.Errorf("поле %s должно быть не меньше %d", fieldName, min)
		}

	case "max":
		max, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return fmt.Errorf("неверный формат max для поля %s: %v", fieldName, err)
		}
		if value > max {
			return fmt.Errorf("поле %s должно быть не больше %d", fieldName, max)
		}
	}

	return nil
}

// Пример использования
type User struct {
	Name  string `validate:"min=3"`
	Age   int    `validate:"min=18;max=65"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

func main() {
	// Тестовые случаи
	fmt.Println("Тест 1 - Короткое имя:")
	if err := Validate(User{Name: "Ив", Age: 18, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	fmt.Println("\nТест 2 - Возраст больше максимума:")
	if err := Validate(User{Name: "Иван", Age: 70, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	fmt.Println("\nТест 3 - Неверный email:")
	if err := Validate(User{Name: "Иван", Age: 35, Email: "invalid email"}); err != nil {
		fmt.Println("Validation error:", err)
	}

	fmt.Println("\nТест 4 - Все поля валидны:")
	if err := Validate(User{Name: "Иван", Age: 35, Email: "test@example.com"}); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Валидация пройдена успешно!")
	}
}
