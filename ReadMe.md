=== Управление СМС через smsfeedback.ru
Реализованы следующие функции</br>
- SendSMS(login, pwd, phone, text, sender, wapurl string)</br>
Отправка СМС.
- GetDelivered(login, pwd string, ids []string)</br>
Проверка доставки.
- GetBalance(login, pwd string) (float32, error)</br>
Получение баланса
