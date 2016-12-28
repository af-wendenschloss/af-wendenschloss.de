# [af-wendenschloss.de](https://af-wendenschloss.de)

Dieses Repository enthält dies Qullen der Seite [af-wendenschloss.de](https://af-wendenschloss.de).

## Pflegehinweise

### Änderungen lokal ansehen

1. In dem Programm IntelliJ IDEA, links im Projekt-Baum die Datei ``local_server.cmd`` suchen.
   Diese sollte sich in der Wurzel des Projektes befinden.
2. Dort mit der rechts Maustaste drauf klicken und dann auf **Run 'local_server.cmd'**
3. Nun sollte sich im unteren Bereich von IntelliJ IDEA ein neuer Bereich öffnen, der ``local_server`` heißt.
   In diesem läuft ein Programm, welches das Projekt lokal zur Verfügung stellt.
4. Nun kann man das Projekt ansehen indem man [localhost:1313](http://localhost:1313/) in seinem
   Browser aufruft. Jede Änderung die man an einer Datei macht, kann man sofort in Browser sehen,
   solange dieses Programm läuft. Einzig im Browser sollte man nach jeder Änderung die Taste ``[F5]``
   drücken, wodurch der Browser die Seite komplett neu lädt.
5. Ist man mit seiner Arbeit fertig klickt man auf das rote Viereck in dem unteren Bereich (``local_server``)
   von IntelliJ IDEA. Hierdurch wird das Skript gestoppt.

### Änderungen veröffentlichen

1. In dem Programm IntelliJ IDEA die Tastenkombination ``[Strg]+[K]`` drücken.
2. In dem folgenden Dialog befindet sich ein Baumdiagramm. Die dort aufgelisteten Dateien
   sollte man kurz überprüfen. Diese stellen alle soeben gemachten Änderungen dar.
     > **Hinweis:** Wählt man eine Datei mit der Maus an und betätigt dann die
     > Tastenkombination ``[Strg]+[D]`` öffnet sich eine separater Dialog in dem man die Änderungen einer
     > Datei genauer inspezieren kann. Mit der Taste ``[Esc]`` schließt sich dieses Fenster wieder.
3. In dem Feld **Commit Message** muss man seine Änderungen kurz beschreiben.
4. Ist alles erledigt, geht man mit der Maus über den Button **Commit** und klickt dann auf den Unterpunkt
   **Commit and push**.
5. In dem folgenden Dialog klickt man dann den Button **Push**.
6. Nun muss man ca. 5 bis 10 Minuten warten, danach sollten die Änderungen auf der Seite zur Verfügung stehen.
