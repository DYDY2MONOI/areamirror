import React from 'react';
import './App.css';

const App = () => {
  return (
    <div className="App">

      <section className="hero">
        <div className="container">
          <h1>Automate your digital world</h1>
          <p>
            Connect your favorite applications and create powerful automations 
            that work for you 24/7. Over 600 connected services.
          </p>
          <div className="hero-buttons">
            <button className="btn-primary">Start for free</button>
            <button className="btn-secondary">See examples</button>
          </div>
        </div>
      </section>

      <section id="features" className="section">
        <div className="container">
          <div className="text-center">
            <h2>How does it work?</h2>
            <p>Create automations in 3 simple steps</p>
          </div>
          <div className="grid grid-3">
            <div className="feature-card">
              <div className="feature-icon">⚡</div>
              <h3>Trigger</h3>
              <p>Choose an event that will trigger your automation. A new email, a social media post, or any other event.</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🔗</div>
              <h3>Connection</h3>
              <p>Connect your favorite services. Gmail, Slack, Twitter, Instagram, and over 600 other applications available.</p>
            </div>
            <div className="feature-card">
              <div className="feature-icon">🎯</div>
              <h3>Action</h3>
              <p>Define what happens next. Send a notification, save a file, or trigger another action.</p>
            </div>
          </div>
        </div>
      </section>

      <section id="services" className="section">
        <div className="container">
          <div className="text-center">
            <h2>Connected services</h2>
            <p>Over 600 applications and services at your disposal</p>
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

      <section id="examples" className="section">
        <div className="container">
          <div className="text-center">
            <h2>Popular applets</h2>
            <p>Discover the most used automations by our community</p>
          </div>
          <div className="grid grid-3">
            <div className="applet-card">
              <div className="applet-title">Automatic backup</div>
              <div className="applet-description">
                Automatically backup all your Instagram photos to Google Drive
              </div>
              <div className="applet-stats">
                <span>🔥 2.3M users</span>
                <span>⭐ 4.8/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Smart reminder</div>
              <div className="applet-description">
                Get a Slack reminder when you receive an important email from your boss
              </div>
              <div className="applet-stats">
                <span>🔥 1.8M users</span>
                <span>⭐ 4.9/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Social sync</div>
              <div className="applet-description">
                Automatically post your tweets to LinkedIn and Facebook
              </div>
              <div className="applet-stats">
                <span>🔥 1.5M users</span>
                <span>⭐ 4.7/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Event management</div>
              <div className="applet-description">
                Automatically create a Google Calendar event from a confirmation email
              </div>
              <div className="applet-stats">
                <span>🔥 1.2M users</span>
                <span>⭐ 4.8/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Weather monitoring</div>
              <div className="applet-description">
                Get a WhatsApp notification if rain is forecast for tomorrow
              </div>
              <div className="applet-stats">
                <span>🔥 950K users</span>
                <span>⭐ 4.6/5</span>
              </div>
            </div>
            <div className="applet-card">
              <div className="applet-title">Sentiment analysis</div>
              <div className="applet-description">
                Automatically analyze the sentiment of your Twitter mentions and receive a report
              </div>
              <div className="applet-stats">
                <span>🔥 800K users</span>
                <span>⭐ 4.5/5</span>
              </div>
            </div>
          </div>
        </div>
      </section>


      <section className="section">
        <div className="container text-center">
          <h2>Ready to automate your life?</h2>
          <p>Join over 2 million users who trust Area</p>
          <div className="hero-buttons">
            <button className="btn-primary">Create free account</button>
            <button className="btn-secondary">See demo</button>
          </div>
        </div>
      </section>

      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-section">
              <h4>Product</h4>
              <ul>
                <li><a href="#features">Features</a></li>
                <li><a href="#integrations">Integrations</a></li>
                <li><a href="#api">API</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Resources</h4>
              <ul>
                <li><a href="#help">Help Center</a></li>
                <li><a href="#tutorials">Tutorials</a></li>
                <li><a href="#blog">Blog</a></li>
                <li><a href="#community">Community</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Company</h4>
              <ul>
                <li><a href="#about">About</a></li>
                <li><a href="#careers">Careers</a></li>
                <li><a href="#press">Press</a></li>
                <li><a href="#contact">Contact</a></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4>Legal</h4>
              <ul>
                <li><a href="#privacy">Privacy</a></li>
                <li><a href="#terms">Terms</a></li>
                <li><a href="#cookies">Cookies</a></li>
                <li><a href="#security">Security</a></li>
              </ul>
            </div>
          </div>
          <div className="footer-bottom">
            <p>&copy; 2024 Area. All rights reserved. Automate your digital world.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;