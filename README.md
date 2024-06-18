#Parallele und verteilte Systeme: Leistungsbeurteilung PVS-01
# Verteiltes Finanzsystem
 
## Übersicht
Dieses Projekt ist Teil der Leistungsbeurteilung für den Kurs "Parallele und verteilte Systeme" im Frühlingssemester 2024. Ziel ist die Entwicklung eines verteilten Finanzsystems, das fiktive Finanzdaten verarbeitet und speichert. Es besteht aus mehreren Docker-Container-basierten Komponenten, die zusammenarbeiten, um eine robuste und skalierbare Lösung zu bieten.
 
## Komponenten
- *Stock Publisher*: Eine Go-Anwendung, die zufällige Finanztransaktionen (Kauf oder Verkauf) für Unternehmen wie Microsoft (MSFT), Tesla (TSLA) und Apple (AAPL) simuliert und diese Daten in RabbitMQ-Queues veröffentlicht.
- *Stock Liveview*: Eine NodeJS-Anwendung, die Aktienpreise in Echtzeit anzeigt. Sie verwendet WebSockets, um die Aktienpreise an den Client zu übertragen, und lädt die Daten aus einem MongoDB-Replica-Set.
- *Consumer*: Eine Anwendung, die Daten aus RabbitMQ liest, Durchschnittspreise berechnet und die aggregierten Daten in MongoDB speichert.
- *MongoDB Cluster*: Ein ausfallsicheres Speichersystem für die aggregierten Finanzdaten.
- *NGINX Load Balancer*: Verteilt die Anfragen auf die Frontend-Instanzen und sorgt für Failover.
 
## Getting Started
Diese Anleitung hilft Ihnen, das Projekt auf Ihrem lokalen Rechner auszuführen.
 
### Prerequisites
Was Sie benötigen, um die Software zu installieren:
- Docker
 
### Installation und Ausführung
1. Klonen Sie die notwendigen Repositories:
    
    git clone https://github.com/kiko46/parallelen.git

 
2. Navigieren Sie in das Verzeichnis mit der docker-compose.yml Datei und starten Sie die Docker-Container:
    
    docker-compose up --build

 
3. Öffnen Sie einen Webbrowser und gehen Sie zu [http://localhost:3000], um die Echtzeit-Aktienpreise zu sehen.
 
## Abschluss
Dieses Projekt demonstriert die Entwicklung und Integration eines verteilten Systems zur Verarbeitung und Anzeige von Finanzdaten. Die Docker-Container ermöglichen eine einfache lokale Ausführung und bieten gleichzeitig Skalierbarkeit und hohe Verfügbarkeit.
