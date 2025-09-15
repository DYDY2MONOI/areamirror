import React from 'react';
import './App.css';

const App = () => {
  return (
    <div className="App">

      {/* Hero Section */}
      <section className="hero">
        <div className="container">
          <h1>Automatisez votre monde numérique</h1>
          <p>
            Connectez vos applications préférées et créez des automatisations puissantes 
            qui travaillent pour vous 24h/24. Plus de 600 services connectés.
          </p>
          <div className="hero-buttons">
            <button className="btn-primary">Commencer gratuitement</button>
            <button className="btn-secondary">Voir les exemples</button>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="section">
        <div className="container">
          <div className="text-center">
            <h2>Comment ça marche ?</h2>
            <p>Créez des automatisations en 3 étapes simples</p>
          </div>
          <div className="grid grid-3">
            <div className="feature-card">
              <div className="feature-icon">⚡</div>
              <h3>Déclencheur</h3>
              <p>Choisissez un événement qui déclenchera votre automatisation. Un nouvel email, un post sur les réseaux sociaux, ou tout autre événement.</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🔗</div>
              <h3>Connexion</h3>
              <p>Connectez vos services préférés. Gmail, Slack, Twitter, Instagram, et plus de 600 autres applications disponibles.</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🎯</div>
              <h3>Action</h3>
              <p>Définissez ce qui se passe ensuite. Envoyer une notification, sauvegarder un fichier, ou déclencher une autre action.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Services Section */}
      <section id="services" className="section">
        <div className="container">
          <div className="text-center">
            <h2>Services connectés</h2>
            <p>Plus de 600 applications et services à votre disposition</p>
          </div>
          <div className="service-grid">
            <div className="service-item">
              <div className="service-icon">📧</div>
              <div className="service-name">Gmail</div>
            </div>
            <div className="service-item">
              <div className="service-icon">💬</div>
              <div className="service-name">Slack</div>
            </div>
            <div className="service-item">
              <div className="service-icon">🐦</div>
              <div className="service-name">Twitter</div>
            </div>
            <div className="service-item">
              <div className="service-icon">📸</div>
              <div className="service-name">Instagram</div>
            </div>
            <div className="service-item">
              <div className="service-icon">📅</div>
              <div className="service-name">Google Calendar</div>
            </div>
            <div className="service-item">
              <div className="service-icon">☁️</div>
              <div className="service-name">Google Drive</div>
            </div>
            <div className="service-item">
              <div className="service-icon">📊</div>
              <div className="service-name">Google Sheets</div>
            </div>
            <div className="service-item">
              <div className="service-icon">🎵</div>
              <div className="service-name">Spotify</div>
            </div>
            <div className="service-item">
              <div className="service-icon">📱</div>
              <div className="service-name">WhatsApp</div>
            </div>
            <div className="service-item">
              <div className="service-icon">💼</div>
              <div className="service-name">LinkedIn</div>
            </div>
            <div className="service-item">
              <div className="service-icon">📺</div>
              <div className="service-name">YouTube</div>
            </div>
            <div className="service-item">
              <div className="service-icon">🔔</div>
              <div className="service-name">Discord</div>
            </div>
          </div>
        </div>
      </section>

      {/* Examples Section */}
      <section id="examples" className="section">
        <div className="container">
          <div className="text-center">
            <h2>Applets populaires</h2>
            <p>Découvrez les automatisations les plus utilisées par notre communauté</p>
          </div>
          <div className="grid grid-3">
            <div className="applet-card">
              <div className="applet-title">Sauvegarde automatique</div>
              <div className="applet-description">
                Sauvegardez automatiquement toutes vos photos Instagram dans Google Drive
              </div>
              <div className="applet-stats">
                <span>🔥 2.3M utilisateurs</span>
                <span>⭐ 4.8/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Rappel intelligent</div>
              <div className="applet-description">
                Recevez un rappel Slack quand vous recevez un email important de votre patron
              </div>
              <div className="applet-stats">
                <span>🔥 1.8M utilisateurs</span>
                <span>⭐ 4.9/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Synchronisation sociale</div>
              <div className="applet-description">
                Publiez automatiquement vos tweets sur LinkedIn et Facebook
              </div>
              <div className="applet-stats">
                <span>🔥 1.5M utilisateurs</span>
                <span>⭐ 4.7/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Gestion d'événements</div>
              <div className="applet-description">
                Créez automatiquement un événement Google Calendar à partir d'un email de confirmation
              </div>
              <div className="applet-stats">
                <span>🔥 1.2M utilisateurs</span>
                <span>⭐ 4.8/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Monitoring météo</div>
              <div className="applet-description">
                Recevez une notification WhatsApp si la pluie est prévue demain
              </div>
              <div className="applet-stats">
                <span>🔥 950K utilisateurs</span>
                <span>⭐ 4.6/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Analyse de sentiment</div>
              <div className="applet-description">
                Analysez automatiquement le sentiment de vos mentions Twitter et recevez un rapport
              </div>
              <div className="applet-stats">
                <span>🔥 800K utilisateurs</span>
                <span>⭐ 4.5/5</span>
              </div>
            </div>
          </div>
        </div>
      </section>


      {/* CTA Section */}
      <section className="section">
        <div className="container text-center">
          <h2>Prêt à automatiser votre vie ?</h2>
          <p>Rejoignez plus de 2 millions d'utilisateurs qui font confiance à Area</p>
          <div className="hero-buttons">
            <button className="btn-primary">Créer un compte gratuit</button>
            <button className="btn-secondary">Voir la démo</button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-section">
              <h4>Produit</h4>
              <ul>
                <li><a href="#features">Fonctionnalités</a></li>
                <li><a href="#integrations">Intégrations</a></li>
                <li><a href="#api">API</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Ressources</h4>
              <ul>
                <li><a href="#help">Centre d'aide</a></li>
                <li><a href="#tutorials">Tutoriels</a></li>
                <li><a href="#blog">Blog</a></li>
                <li><a href="#community">Communauté</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Entreprise</h4>
              <ul>
                <li><a href="#about">À propos</a></li>
                <li><a href="#careers">Carrières</a></li>
                <li><a href="#press">Presse</a></li>
                <li><a href="#contact">Contact</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Légal</h4>
              <ul>
                <li><a href="#privacy">Confidentialité</a></li>
                <li><a href="#terms">Conditions</a></li>
                <li><a href="#cookies">Cookies</a></li>
                <li><a href="#security">Sécurité</a></li>
              </ul>
            </div>
          </div>
          <div className="footer-bottom">
            <p>&copy; 2024 Area. Tous droits réservés. Automatisez votre monde numérique.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;