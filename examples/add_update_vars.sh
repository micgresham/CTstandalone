#!/bin/bash
filename=$1
variable="site_id"
while read line; do

  # reading each line and get the serials
  ../bin/CTsite_serials-linux --site=$line --dev_type=ALL

  # process each site serial file and add/update the requested vars
  ../bin/CTadd_var-linux --infile $line-serials.txt --variable=$variable --value=$line --test

	
done < $filename
