Description

Ce projet Python génère des fichiers de configuration de démarrage pour 10 routeurs répartis dans différents Systèmes Autonomes (AS). 
Il permet de configurer des topologies complexes avec des protocoles de routage et des sessions BGP.

Nous avons comme fonctionnalités : 
    -> La mise en place d'addresse IP'.
    -> La mise en place de plusieurs AS (AS_Y,AS_X,AS_Z,AS_Z2) and leur numéro.
    -> La mise en place des protocoles de routage OSPF (Pour AS_Y) et RIP (Pour AS_X).
    -> Nous avons mis iBGP dans AS 
    -> Nous avons établi des sessions eBGP entre deux les deux AS (AS_X et AS_Y)
    -> Automatisation d'interface loopback'