# not-go-github
 Библиотека для управления GitHub на Go, использующая [клиентсую библиотеку для взаимодействия с GitHub API v3](github.com/google/go-github).

## Пример использования
Для использования функций, получим GitServiceIFace:
```go
ghs, err := NewGitHubService(context.TODO())
```
Теперь можем использовать реализованные методы:

```go
// Получить список соавторов репозитория "google/go-github"
contributors, err := ghs.GetRepositoryContributors("google", "go-github")
```
```go
// Получить информацию о репозитории "jostanise/rsa_encrypted_local_chat"
repo, err := ghs.GetRepositoryByName("jostanise", "rsa_encrypted_local_chat")
```


# Тестирование ограничений доступа к GitHub API

### Тестирование с токеном
```go
// Используем токен доступа, поместив его в файл .env
godotenv.Load(".env")
```
Тестируем:
```go
ghs, _ := NewGitHubService(context.Background())

for i := 0; ; i++ {
	_, err := ghs.GetUserRepositories("jostanise")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
```
Тестирование было преждевременно остановлено на основании промежуточного  результата:

![Результат тестирования с токеном доступа](https://user-images.githubusercontent.com/82827415/179687521-025cb2c6-2794-489e-9373-a7290b5ac5d9.png)

### Тестирование без токена
Промежуточный результат:

![Результат тестирования без токена доступа](https://user-images.githubusercontent.com/82827415/179687608-91fe0383-1688-43e5-8b59-bbd95f3f5b75.png)

## Вывод
При отсутствии токена доступа, количество запросов, которые можно произвести за час, вместо `5000` ограничивается до `60`, а также запросы возвращают ошибку. После достижения ограничения по количеству запросов на один IP-адрес, запросы возвращают ошибки.
