# Opis projekta:
## Naslov: 
Wikno notes

## Člani skupine in številka projektne skupine
Člani: Blaž Bergant 
Številka: 06

## Povezava do GitHub organizacije in repozitorijev
https://github.com/BwezB/Wikno-notes

## Kratek opis projekta
Wikno notes je inovativna aplikacija za učenje in organizacijo, ki vsem uporabnikom omogoča urejanje svojih zapiskov o skupnih entitetah in kreiranje povezav med entitetami. Vsak posameznik bo tako lahko kreiral svojo mrežo znanja (entitet z nekimi lastnostmi, ki se med seboj povezujejo), naenkrat pa si bo lahko pogledal zapiske drugih uporabnikov da razširi svoje znanje. V kontekstu organizacije, to pomeni da se bo lahko organiziral z todo listi za zasebne projekte, naenkrat pa bo projekte in todo liste lahko delil z drugimi uporabniki. 
Cilj projekta je narediti platformo, ki bo revolucionirala kako ljudje razmišljamo o zapiskih, učenju in organizaciji.

## Ogrodje in razvojno okolje
- Backend: Golang -> Jezik je enostaven za učenje in vse bolj popularen za backende.
- Frontend: Swift -> Z jezikom hitro zgradimo estetske aplikacije ki delujejo dobro. 

## Shema arhitekture
![alt text](images/image.png)

## Seznam funkcionalnosti mikrostoritev

1. Entity, Property Manager: 
    - CRUD za entitete, povezave, lastnosti.
    - Ugotavlja duplikate v podatkih, entitetah.
    - Pošilja in prejema vse podatke ki gredo v podatkovno bazo in iz nje

2. Natural language processing Manager
    - Čekira duplikate pri novo ustvarhenih tipih povezav in tipih lastnosti (če že obstaja tip z enakim pomenom)
    - Išče kateri uporabniki delajo podobne zapiske, da jim privzeto kaže zapiske drug od drugega (ko odkrivajo novo znanje)

3. Retention Manager:
    - Predlaga flashcarde za uporabnika (glede na že pridobljeno znanje)
    - Predlaga vaje za uporabnika (glede na že pridobljeno znanje)

4. Users Manager:
    - CRUD za uporabnike
    - Avtentikacija
    - Ali ima user pravico dobiti neke informacije


## Primeri uporabe

### Kratki
1. Uporabnik ustvari javno entiteto "Turingov stroj" in ji doda lastnost tipa "Opis" z vsebino "Je model računanja..."
2. Uporabnik entiteti "Turingov stroj" doda povezavo tipa "Kreator" ki kaže na entiteto "Alan Turing"
3. Uporabnik ustvari zasebno entiteto "Pomij posodo", z lastnostmi: Entity type: Task; Priority: A; Deadline: Today. Da se opomni o tem opravilu.
4. Uporabnik prebere lastnosti (npr. "Opis") o entiteti "Alan Turing" od drugega uporabnika ki dela podobne zapiske kot on. 
5. Uporabnik ustvari nov tip povezave "Tata" z opisom "Starš moški". Aplikacija ga vpraša, če je to enaka entiteta kot "Oče".


### Kompkensni primer:

**Uporabniki:**
- Ana: Ravno začela z predmetom "Algoritmi in podatkovne strukture"
- Bojan Predmet "Algoritmi in podatkovne strukture" že opravil
- Domen Se pripravlja na 4 rok izpita.
- Cvetka Ravno začela z predmetom.

1. **Ustvarjanje in povezovanje entitet:**
   Ana ustvari javno entiteto "Algoritmi in podatkovne strukture" ki ji doda lastnosti in povezave:
   - povezava "Tip": "Predmet"
   - povezava "Podpirajoča uztanova": "FRI"
   - lastnost "Opis": "Predmet ki obdela osnovne algoritme in podatkovne strukture"
   - lastnost "Semester": "3"

2. **Sodelovanje in deljenje znanja:**
   Bojan, ki je predmet že opravil, z napiše "Predmet se gre o časovni zahtevnosti, prostorski zahtevnosti, drevesih." in s tem kreira povezave:
   - povezava "Se gre o": "Časovna zahtevnost"
   - povezava "Se gre o": "Prostorska zahtevnost"

4. **Pregled in razširitev znanja:**
   Domen, ki se pripravlja na izpit doda entitete povezane z ključnimi koncepti. Kreira "Rdeče črno drevo". Doda mu lastnosti:
    - lastnost "Kako deluje": "..."
    - lastnost "Uporaba": "..."
    - povezava "Je implementacija od": "Drevesa"

6. **Uporaba NLP za odkrivanje podobnosti:**
    Cvetka je nevešča, in kreira javno entiteto "APS", z lastnostmi in povezavami:
    - povezava "Tip": "Predmet"
    - povezava "Podpirajoča uztanova": "FRI"
    - lastnost "Semester": "3"
    Ker ima entiteta podobne povezave kot že znana entiteta "Algoritmi in podatkovne strukture", "NLP Manager" predlaga da ju združi.

7. **Uporaba retention managerja:**
    Ana, ki se pripravlja na izpit naredi par entitet z povezavo tipa "Tip": "Flashcard", ki imajo lastnost "Vprašanje" in "Odgovor". Retention manager ji predlaga dodatne entitete tipa "Flashcard" ki so jih kreirali drugi uporabniki.
