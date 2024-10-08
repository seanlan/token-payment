dbuser="root"
dbpwd="q145145145"
dbhost="127.0.0.1:3306"
dbname="tokenpay"
conn="mysql://$dbuser:$dbpwd@tcp($dbhost)/$dbname?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci"
package=$(go list -m)
prefix="t_"
workdir=$(dirname $0)
template=$workdir"/templates"
modelPackage="sqlmodel"
modelPath="internal/dao/sqlmodel"
daoPackage="dao"
daoPath="internal/dao"
lazy dao --conn=$conn --database=$dbname --prefix=$prefix --package=$package --template=$template \
 --model=$modelPackage --model-path=$modelPath --dao=$daoPackage --dao-path=$daoPath
