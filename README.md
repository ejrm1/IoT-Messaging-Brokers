# Tarea 2 Sistemas Distribuidos, IoT-Messaging-Brokers

Eduardo Riquelme Molina
Manuel Caceres Farias


[Video](https://youtu.be/JcGpvW6ct-U)
# Utilizacion
Para correr los docker compose, simplemente dirijase a la carpeta del compose que desea ejecutar(Kafka, RabbitMQ, o con Types) y ejecute el siguiente comando:
``` docker-compose build up --build```
El programa se ejecutara automaticamente

Es importante notar, que en caso de querer ejecutarse en Windows, hay que hacer una leve modificacion a los dockerfiles de RabbitMQ:
```
#RUN sed -i 's/\r$//' ./wait-for-it.sh  && \  
#        chmod +x ./wait-for-it.sh

#RUN sed -i 's/\r$//' ./start.sh  && \  
#        chmod +x ./start.sh
```
Estas lineas de codigo deben ser descomentadas en tanto el dockerfile de productor como el del consumidor.
