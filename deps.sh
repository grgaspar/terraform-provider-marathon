missing_package="start"
while [ "$missing_package" != "" ]
do
  missing_package=$(godep save 2>&1 | \
    egrep '^godep: Package (.*) not found' | \
    sed 's/.*(\(.*\)).*/\1/'); 
 [ "$missing_package" != "" ] && { 
   echo "Installing missing package: ${missing_package}" ; 
   go get -u "${missing_package}"
  }
done
godep save
