go mod init appointment-system
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql



db connection
SELECT user, host, plugin FROM mysql.user WHERE user = 'root';
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'P@mmount847';
FLUSH PRIVILEGES;
