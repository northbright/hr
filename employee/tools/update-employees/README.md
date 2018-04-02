# update-employees

update-employees is a program which read and parse employee data in given CSV file and update the data stored in Redis database. It's written in [Golang](http://golang.org).

#### Convention of Employee CSV File
* CSV files should have head row.
* First 4 columns(0-based):
  * column 0: name of employee
  * column 1: sex of employee
  * column 2: ID card number of employee
  * column 3: mobile phone number
* Example:

| 姓名 | 性别 | 身份证号码 | 手机号码 |
| :--- | :--- | :--- | :--- |
| 张三 | 男 | 310104198101010000 | 13500000000 |
| 王强 | 男 | 310104198201010000 | 13600000000 |
| 小红 | 女 | 310104198301010000 | 13700000000 |
| 小明 | 男 | 310104198401010000 | 13800000000 |

#### Usage

    -a string
    	Redis server address. Ex: -a='127.0.0.1:6379'
    -f string
    	path of CSV file which contains employee data. Ex: -f='employees.csv'
    -p string
    	Redis password. Ex: -p='my_password'

    usage:
        update-employees -a=<Redis server address> -p=<Redis password> -f=<csv file>

#### Example

* [`files/original.csv`](files/original.csv) contains the original employee data.
* [`files/update.csv`](files/update.csv) contains the updated employee data.
  * Change `张三`'s mobile phone number from `13500000000` to `13900000000`.

* Commands

          update-employees -a='127.0.0.1:6379' -p='' -f='files/files/original.csv'
          update-employees -a='127.0.0.1:6379' -p='' -f='files/update.csv'




