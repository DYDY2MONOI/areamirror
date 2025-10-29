<template>
  <footer class="site-footer">
    <div class="footer-content">
      <div class="footer-section">
        <div class="footer-logo">
          <h3 class="company-name">{{ companyName }}</h3>
          <p class="company-tagline">{{ companyTagline }}</p>
        </div>
        <p class="footer-description">
          {{ companyDescription }}
        </p>
      </div>

      <div v-for="section in footerSections" :key="section.title" class="footer-section">
        <h4 class="footer-title">{{ section.title }}</h4>
        <ul class="footer-links">
          <li v-for="link in section.links" :key="link.text">
            <a href="#" @click.prevent="$emit('link-click', link.id)">{{ link.text }}</a>
          </li>
        </ul>
      </div>

      <div class="footer-section">
        <h4 class="footer-title">{{ contactTitle }}</h4>
        <div class="contact-info">
          <div v-for="contact in contactInfo" :key="contact.type" class="contact-item">
            <v-icon size="16" color="var(--color-accent-primary)">{{ contact.icon }}</v-icon>
            <span>{{ contact.value }}</span>
          </div>
        </div>
        <div class="social-links">
          <a
            v-for="social in socialLinks"
            :key="social.id"
            href="#"
            class="social-link"
            @click.prevent="$emit('social-click', social.id)"
          >
            <v-icon size="20">{{ social.icon }}</v-icon>
          </a>
        </div>
      </div>
    </div>

    <div class="footer-bottom">
      <div class="footer-bottom-content">
        <p class="copyright">
          © {{ currentYear }} {{ companyName }}. All rights reserved.
        </p>
        <div class="footer-bottom-links">
          <a
            v-for="link in bottomLinks"
            :key="link.id"
            href="#"
            @click.prevent="$emit('bottom-link-click', link.id)"
          >
            {{ link.text }}
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
interface FooterLink {
  id: string
  text: string
}

interface FooterSection {
  title: string
  links: FooterLink[]
}

interface ContactInfo {
  type: string
  icon: string
  value: string
}

interface SocialLink {
  id: string
  icon: string
}

interface Props {
  companyName?: string
  companyTagline?: string
  companyDescription?: string
  contactTitle?: string
  footerSections?: FooterSection[]
  contactInfo?: ContactInfo[]
  socialLinks?: SocialLink[]
  bottomLinks?: FooterLink[]
}

interface Emits {
  (e: 'link-click', linkId: string): void
  (e: 'social-click', socialId: string): void
  (e: 'bottom-link-click', linkId: string): void
}

const currentYear = new Date().getFullYear()

withDefaults(defineProps<Props>(), {
  companyName: 'AREA',
  companyTagline: 'Intelligent Automation Platform',
  companyDescription: 'Connect your favorite services with intelligent automation. Build powerful workflows that work for you.',
  contactTitle: 'Contact',
  footerSections: () => [
    {
      title: 'Product',
      links: [
        { id: 'features', text: 'Features' },
        { id: 'integrations', text: 'Integrations' },
        { id: 'api', text: 'API' },
        { id: 'documentation', text: 'Documentation' }
      ]
    },
    {
      title: 'Company',
      links: [
        { id: 'about', text: 'About Us' },
        { id: 'careers', text: 'Careers' },
        { id: 'blog', text: 'Blog' },
        { id: 'press', text: 'Press' }
      ]
    },
    {
      title: 'Support',
      links: [
        { id: 'help', text: 'Help Center' },
        { id: 'community', text: 'Community' },
        { id: 'contact', text: 'Contact' },
        { id: 'status', text: 'Status' }
      ]
    }
  ],
  contactInfo: () => [
    { type: 'email', icon: 'mdi-email', value: 'contact@area.com' },
    { type: 'phone', icon: 'mdi-phone', value: '+33 7 41 61 72 18' },
    { type: 'location', icon: 'mdi-map-marker', value: 'Paname, France' }
  ],
  socialLinks: () => [
    { id: 'twitter', icon: 'mdi-twitter' },
    { id: 'github', icon: 'mdi-github' },
    { id: 'linkedin', icon: 'mdi-linkedin' },
    { id: 'discord', icon: 'mdi-discord' }
  ],
  bottomLinks: () => [
    { id: 'privacy', text: 'Privacy Policy' },
    { id: 'terms', text: 'Terms of Service' },
    { id: 'cookies', text: 'Cookie Policy' }
  ]
})

defineEmits<Emits>()
</script>

<style scoped>
.site-footer {
  background: var(--color-bg-secondary);
  border-top: 1px solid var(--color-border-primary);
  margin-top: 4rem;
  padding: 3rem 0 0 0;
  backdrop-filter: blur(20px);
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1.5fr;
  gap: 3rem;
  margin-bottom: 3rem;
  padding: 0 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.footer-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.footer-logo {
  margin-bottom: 1rem;
}

.company-name {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.5rem 0;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 50%, #ec4899 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  filter: drop-shadow(0 2px 8px rgba(59, 130, 246, 0.3));
}

.company-tagline {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 500;
}

.footer-description {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.footer-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 1rem 0;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.footer-links li a {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.footer-links li a:hover {
  color: var(--color-accent-primary);
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.social-links {
  display: flex;
  gap: 1rem;
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: var(--transition-normal);
}

.social-link:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
  color: var(--color-accent-primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow);
}

.footer-bottom {
  border-top: 1px solid var(--color-border-primary);
  padding: 2rem 0;
}

.footer-bottom-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  padding: 0 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.copyright {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.footer-bottom-links {
  display: flex;
  gap: 2rem;
}

.footer-bottom-links a {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.footer-bottom-links a:hover {
  color: var(--color-accent-primary);
}

@media (max-width: 1024px) {
  .footer-content {
    grid-template-columns: 1fr 1fr 1fr;
    gap: 2rem;
  }

  .footer-section:first-child {
    grid-column: 1 / -1;
  }
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    gap: 2rem;
    padding: 0 1rem;
  }

  .footer-bottom-content {
    flex-direction: column;
    text-align: center;
    padding: 0 1rem;
  }

  .footer-bottom-links {
    justify-content: center;
  }

  .social-links {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .site-footer {
    padding: 2rem 0 0 0;
  }

  .footer-content {
    gap: 1.5rem;
    padding: 0 1rem;
  }

  .footer-bottom {
    padding: 1.5rem 0;
  }

  .footer-bottom-content {
    padding: 0 1rem;
  }

  .footer-bottom-links {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
