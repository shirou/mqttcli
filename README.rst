mqttcli -- MQTT Client for shell scripting
=================================================

mqttcli is an MQTT 3.1.1 Client which has almost same options with
mosquitto_pub/sub. However, it has additional functionallity and a
pubsub command which is suite for the shell script pipelining.

Install
==============

Download from here. Please choose your artitecture. (and chmod ugo+x if needed)

https://drone.io/github.com/shirou/mqttcli/files

Or if you have golang environment.

::

  go get github.com/shirou/mqttcli


Usage
==============

common
----------

You can set host, port, username and password on the Environment.

::

    export MQTT_HOST="localhost"
    export MQTT_PORT="1883"
    export MQTT_USERNAME="user"
    export MQTT_PASSWORD="blahblah"

or using this config file. You can specifiy config file using
``--conf`` option, or place ``~/.mqttcli.cfg``.

::

  % mqttcli sub --conf settings.json  -t "some/topic"

example:

::

   {
     "host": "localhost",
     "port": 1883,
     "username": "user",
     "password": "blahblah"
   }


Pub
-------

::

  mqttcli pub -t "some/where" -m "your message"

  or

  tail -f /var/log/nginx.log | mqttcli pub -t "some/where" -s

`-s` is diffrent from mosquitto_pub, it sends one line to one message.

Sub
------

::

  mqttcli sub -t "some/#"


PubSub
---------

Note: This subcommand is just a concept work. Might be delete in the future.

Publish from stdin AND Subscribe from some topics and print stdout.

::

  tail -f /vag/log/nginx.log | mqttcli pubsub --pub "some/a" --sub "some/#" > filterd.log

This is useful when other client manuplate something and send back to
the topic.


Reference
==============

paho.mqtt.golang.git
  http://godoc.org/git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git


License
===========

Eclipse Public License - v 1.0 (same as Paho's)

