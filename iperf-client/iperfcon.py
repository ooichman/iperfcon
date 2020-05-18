from http.server import HTTPServer, BaseHTTPRequestHandler

import os
import re
import subprocess
import json
import traceback

# supported output formats and their HTTP Content-Type
contentTypes = {'plain': 'text/plain', 'html': 'text/html', 'json': 'application/json'}
binaryPrefixes = {'k':'K-bit', 'm':'M-bit', 'g':'G-bit', 't':'T-bit'}

class IperfCon:
    output = None
    iperfServer = None
    iperfPort = None
    warningValue = None
    criticalValue = None
    iperfFormat = None
    acceptFormat = None

    def __init__(self, cmd, params, acceptFormat, queryString):
        self.cmd = cmd
        self.params = params
        self.acceptFormat = acceptFormat
        # next item is for backwards compatibility
        self.queryString = queryString

    @staticmethod
    def _validateBinaryPrefix(arg, value):
        if value not in binaryPrefixes:
            return 'Illegal binary prefix for "' + arg + '" value: ' + \
                value + '. Supported values are: ' + str(binaryPrefixes)
        return None

    def _initialize(self):
        if self.cmd.find('/api.cgi') < 0 and self.cmd.find('/state') < 0:
            return 404, None, None

        self.iperfServer = self.params.get('server')
        self.iperfPort = self.params.get('port')
        if not self.iperfServer or not self.iperfPort:
            return 422, 'Missing query parameter(s): ' + \
                ('server' if not self.iperfServer else '') + \
                (' port' if not self.iperfPort else ''), contentTypes.get('plain')

        formatType = self.params.get('type')
        if formatType:
            self.acceptFormat = formatType
        if self.acceptFormat.find('plain') < 0 and \
                self.acceptFormat.find('html') < 0 and \
                self.acceptFormat.find('json') < 0:
            return 422, 'Unsupported content-type: ' + self.acceptFormat, \
                    contentTypes.get('plain')

        if self.cmd.find('/state') >= 0:
            warningQuery = self.params.get('warning')
            criticalQuery = self.params.get('critical')
            if not warningQuery or not criticalQuery:
                return 422, 'Missing query parameter(s): ' + \
                    ('warning' if not warningQuery else '') + \
                    (' critical' if not criticalQuery else ''), \
                    contentTypes.get('plain')
            warningBinaryPrefix = warningQuery[-1]
            criticalBinaryPrefix = criticalQuery[-1]

            content = self._validateBinaryPrefix('warning', warningBinaryPrefix)
            if content:
                return 422, content, contentTypes.get('plain')
            content = self._validateBinaryPrefix('critical', criticalBinaryPrefix)
            if content:
                return 422, content, contentTypes.get('plain')
            if warningBinaryPrefix != criticalBinaryPrefix:
                return 422, 'warning/critical binary prefixes must match. Values provide: ' + \
                    warningBinaryPrefix + ' and ' + criticalBinaryPrefix, contentTypes.get('plain')
            self.iperfFormat = warningBinaryPrefix

            try:
                self.warningValue = int(warningQuery[:-1])
            except ValueError:
                return 422, 'warning value "' + warningQuery[:-1] + '" is not a valid number', \
                    contentTypes.get('plain')
            try:
                self.criticalValue = int(criticalQuery[:-1])
            except ValueError:
                return 422, 'critical value "' + criticalQuery[:-1] + '" is not a valid number', \
                    contentTypes.get('plain')
            if self.criticalValue >= self.warningValue:
                return 422, 'critical value must be smaller than warning value', \
                    contentTypes.get('plain')

        return 0, None, None

    def getResults(self):
        statusCode, content, contentFormat = self._initialize()
        if statusCode != 0:
            return statusCode, content, contentFormat

        statusCode, content, contentFormat = self._runCommand()
        if statusCode != 0:
            return statusCode, content, contentFormat

        self.output = content

        if self.cmd.find('/api.cgi') >= 0:
            # backwards compatible
            statusCode, content, contentFormat = self._getResultsApiCgi()
        elif self.cmd.find('/state') >= 0:
            statusCode, content, contentFormat = self._getResultsState()
        return statusCode, content, contentFormat

    # def getOutput(self):
    #     return self.output

    def _runCommand(self):
        cmd = ['iperf3', '-c', self.iperfServer, '-p', self.iperfPort]
        if self.iperfFormat:
            cmd.append('-f')
            cmd.append(self.iperfFormat)
        try:
            cmdOutput = subprocess.check_output(cmd, stderr=subprocess.STDOUT).decode()
            return 0, cmdOutput, ''
        except subprocess.CalledProcessError as error:
            message = error.output.decode()
        except:
            message = 'Unable to run iperf3'
        return 503, message, contentTypes.get('plain')

    def _getSummaryLines(self):
        idFlag = 0
        summaryLines = []
        for line in self.output.split('\n'):
            if line.find('[ ID]') >= 0:
                idFlag += 1
            if idFlag >= 2 and idFlag < 5:
                summaryLines.append(line)
                idFlag += 1
        return summaryLines

    def _formatHtml(self):
        htmlString = '<!DOCTYPE html>' + \
	        '<html><head>' + \
    	    '<title>Bash-CGI Example 1</title>' + \
    	    '</head><body>' + \
            '<h1>The iperf output</h1>' + \
    	    '<p>QUERY_STRING: ' + self.queryString + \
            '<br>var1=' + self.iperfServer + \
            '<br>var2=' + self.iperfPort + '</p>'
        for line in self._getSummaryLines():
            htmlString += '<p>' + line + '</p>'
        htmlString += '<hr>' + '</body></html>'
        return htmlString

    def _formatJson(self):
        outputDict = {}
        regex = re.compile('[ \t]+')
        resultCount = 0
        for line in self._getSummaryLines():
            elements = regex.split(line)
            if len(elements) >= 9:
                lineDict = {}
                lineDict['id'] = elements[1].replace(']', '')
                lineDict['Interval'] = elements[2]
                lineDict['Transfer'] = elements[4] + ' ' + elements[5]
                lineDict['Bandwidth'] = elements[6] + ' ' + elements[7]
                lineDict['Host'] = elements[len(elements) - 1]
                if len(elements) == 10:
                    lineDict['Retr'] = elements[8]
                outputDict['result' + str(resultCount)] = lineDict
                resultCount += 1
        return json.dumps(outputDict)

    def _getResultsApiCgi(self):
        if self.acceptFormat.find("html") >= 0:
            return 200, self._formatHtml(), contentTypes.get('html')
        if self.acceptFormat.find("json") >= 0:
            return 200, self._formatJson(), contentTypes.get('json')
        return 200, '\n'.join(self._getSummaryLines()), contentTypes.get('plain')

    def _getResultsState(self):
        status = None
        units = None
        results = {}
        regex = re.compile('[ \t]+')
        # get bitrates for both the sender and receiver
        for line in self._getSummaryLines():
            elements = regex.split(line)
            if len(elements) >= 9:
                units = elements[7]
                results[elements[len(elements) - 1]] = float(elements[6])
        for bitrate in results.values():
            if bitrate < self.criticalValue:
                status = 'CRITICAL'
                break
        if not status:
            for bitrate in results.values():
                if bitrate < self.warningValue:
                    status = 'WARNING'
        if not status:
            status = 'OK'
        if self.acceptFormat.find("json") >= 0:
            results['units'] = units
            results['status'] = status
            return 200, json.dumps(results), contentTypes.get('json')
        return 200, status, contentTypes.get('plain')

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):

    def do_GET(self):
        queryComponents = re.split(',|\?', self.path)
        cmd = queryComponents.pop(0)
        try:
            params = dict((k.lower(), v) for k,v in (x.split("=") for x in queryComponents))
        except:
            self.send_response(422)
            return

        acceptFormat = self.headers.get('Accept')
        if acceptFormat == '*/*':
            acceptFormat = contentTypes.get('plain')

        iperfCon = IperfCon(cmd, params, acceptFormat, \
                'http://' + self.headers.get('Host') + self.path)

        try:
            statusCode, content, contentFormat = iperfCon.getResults()
        except Exception as ex:
            # safety net for uncaught exceptions
            statusCode = 503
            content = ''.join(traceback.format_exception(etype=type(ex), \
                value=ex, tb=ex.__traceback__))
            contentFormat = contentTypes.get('plain')

        self.send_response(statusCode)
        if contentFormat:
            self.send_header('Content-Type', contentFormat)
        self.end_headers()
        if content:
            self.wfile.write((content + '\n').encode())
        return

if 'IPERF_PORT' in os.environ:
    httpPort = int(os.environ['IPERF_PORT'])
else:
    httpPort = 8080
print('starting HTTP server on port ' + str(httpPort))
httpDaemon = HTTPServer(('', httpPort), SimpleHTTPRequestHandler)
httpDaemon.serve_forever()
