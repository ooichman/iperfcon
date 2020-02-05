#!/bin/bash
    
# CORS is the way to communicate, so lets response to the server first
echo "Content-type: text/html"    # set the data-type we want to use
echo ""    # we dont need more rules, the empty line initiate this.

# CORS are set in stone and any communication from now on will be like reading a html-document.
# Therefor we need to create any stdout in html format!
    
# create html scructure and send it to stdout
# The content will be created depending on the Request Method 
if [ "$REQUEST_METHOD" = "GET" ]; then
   
    # Note that the environment variables $REQUEST_METHOD and $QUERY_STRING can be processed by the shell directly. 
    # One must filter the input to avoid cross site scripting.
	if [[ ! $QUERY_STRING =~ 'server' ]]; then
		echo "No Service IP address has been given"
	elif [[ ! $QUERY_STRING =~ 'port' ]]; then     
		echo "No Service port has been given"
	fi
    	IFS=',' read -r -a array <<< "$QUERY_STRING"

    for element in "${array[@]}"; do
	if [[ $element =~ 'server' ]]; then
	Server=$(echo "$element" | awk -F \= '{print $2}')  # read value of "var1"
	fi
#	Server_Dec=$(echo -e $(echo "$Server" | sed 's/+/ /g;s/%\(..\)/\\x\1/g;'))    # html decode

	if [[ $element =~ 'port' ]]; then
	Port=$(echo "$element" | awk -F \= '{print $2}')
	fi

	if [[ $element =~ 'type' ]]; then
		Type=$(echo "$element" | awk -F \= '{print $2}')
	else
		Type="html";
	fi
    done
    # create content for stdout
    if [ ${Type} == "html" ]; then
	echo "<!DOCTYPE html>"
	echo "<html><head>"
    	echo "<title>Bash-CGI Example 1</title>"
    	echo "</head><body>"
    	echo "<h1>The iperf output</h1>"
    	echo "<p>QUERY_STRING: ${QUERY_STRING}<br>var1=${Server}<br>var2=${Port}</p>"    # print the values to stdout
  
    	flag=0
    	/usr/bin/iperf3 -c ${Server} -p ${Port} | while read LINE; do
    		if [[ $LINE =~ 'ID' ]]; then
			flag=$(($flag+1))
    		fi
    		if [ $flag -eq 2 ]; then
    			echo "<p>$LINE</p>"
    		fi
    	done 
    fi

    if [ $Type == "json" ]; then 
	echo "{"
	linenum=0
	/usr/bin/iperf3 -c ${Server} -p ${Port} |  egrep -v "ID|iperf" | grep -v "^$" | tail -2 | while read LINE; do
		echo "    \"result${linenum}\": {"
		linenum=$(($linenum+1))
		ID=$(echo $LINE | awk -F \[ '{print $2}' | awk -F \] '{print $1}'| awk '{print $1}'); 
		INTERVAL=$(echo $LINE | awk '{print $3}');
      		TRANSFER=$(echo $LINE | awk '{print $5" "$6}');
      		BANDWIDTH=$(echo $LINE | awk '{print $7" "$8}')
      		RETR=$(echo $LINE | awk '{print $9" "$10}')
	echo "\"id\": \"${ID}\",  \"Interval\": \"${INTERVAL}\", \"Transfer\": \"${TRANSFER}\", \"Bandwidth\": \"${BANDWIDTH}\", \"Retr\": \"${RETR}\""
		if [[ $linenum -eq 1 ]]; then
			echo "},"
		fi
	done
        echo "    }"
	echo "}"
    fi
else

    echo "<title>456 Wrong Request Method</title>"
    echo "</head><body>"
    echo "<h1>456</h1>"
    echo "<p>Requesting data went wrong.<br>The Request method has to be \"GET\" only!</p>"
fi

if [ ${Type} == "html" ]; then
	echo "<hr>"
	echo "$SERVER_SIGNATURE"    # an other environment variable
	echo "</body></html>"    # close html
fi
    
exit 0
