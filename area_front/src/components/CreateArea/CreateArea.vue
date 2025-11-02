<template>
  <div class="create-area-container">
    <div class="main-card">
      <div class="card-header">
        <div class="header-content">
          <div class="header-text">
            <h1 class="card-title">Create New Area</h1>
            <p class="card-subtitle">Connect your favorite services with intelligent automation</p>
          </div>
          <div class="header-actions">
            <button class="info-button" @click="showGuide = true" title="Show Guide">
              <v-icon size="20" color="white">mdi-information</v-icon>
            </button>
            <button class="close-button" @click="$emit('close')">
              <v-icon size="24" color="white">mdi-close</v-icon>
            </button>
          </div>
        </div>
      </div>
      <div class="card-content">
        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-tag-outline</v-icon>
            <span class="label-text">Area Details</span>
          </div>

          <div class="input-group">
            <div class="input-container">
              <label class="input-label">Area Name</label>
              <input
                v-model="form.areaName"
                class="modern-input"
                placeholder="Enter a name for your area"
                required
              />
            </div>

            <div class="input-container">
              <label class="input-label">Description</label>
              <textarea
                v-model="form.description"
                class="modern-textarea"
                placeholder="Describe what this area does..."
                rows="3"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-link-variant</v-icon>
            <span class="label-text">Service Connection</span>
          </div>

          <div class="connection-flow">
            <div class="service-selector">
              <div class="selector-header">
                <div class="selector-icon">
                  <v-icon size="24" color="#3b82f6">mdi-play-circle</v-icon>
                </div>
                <div class="selector-info">
                  <h3 class="selector-title">When this happens</h3>
                  <p class="selector-subtitle">Choose the trigger service</p>
                </div>
              </div>

              <div class="service-selection">
                <div v-if="servicesError" class="service-error">
                  <v-icon size="18" color="#ef4444">mdi-alert-circle</v-icon>
                  <span>{{ servicesError }}</span>
                </div>

                <div v-else-if="isLoadingServices" class="service-loading">
                  <div class="spinner"></div>
                  <span>Loading services...</span>
                </div>

                <div v-else>
                  <div v-if="!form.triggerService" class="service-grid">
                    <div
                      v-for="item in triggerItems.slice(0, 15)"
                      :key="item.value"
                      class="service-card"
                      @click="selectTrigger(item.value)"
                    >
                      <div class="service-card-icon">
                        <img :src="item.icon" :alt="item.title" class="service-icon" />
                      </div>
                      <span class="service-card-name">{{ item.title }}</span>
                    </div>
                    <div class="service-card more-services" @click="showAllTriggerServices = true" v-if="triggerItems.length > 15">
                      <div class="service-card-icon">
                        <v-icon size="24" color="#3b82f6">mdi-plus</v-icon>
                      </div>
                      <span class="service-card-name">More...</span>
                    </div>
                  </div>

                  <div v-else class="selected-service-display">
                    <div class="selected-service-card">
                      <div class="service-avatar">
                        <img :src="serviceIcons[form.triggerService || ''] || getFallbackIcon(form.triggerService || '')" :alt="getServiceName(form.triggerService)" class="service-icon" />
                      </div>
                      <div class="service-info">
                        <span class="service-name">{{ getServiceName(form.triggerService) }}</span>
                        <span class="service-type">Trigger Service</span>
                      </div>
                      <button class="change-service-btn" @click="form.triggerService = ''">
                        <v-icon size="16">mdi-close</v-icon>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="connection-arrow">
              <div class="arrow-line"></div>
              <div class="arrow-icon">
                <v-icon size="20" color="#3b82f6">mdi-arrow-down</v-icon>
              </div>
              <div class="arrow-line"></div>
            </div>

            <div v-if="form.triggerService" class="service-selector intermediate-selector">
              <div class="selector-header">
                <div class="selector-icon">
                  <v-icon size="24" color="#8b5cf6">mdi-brain</v-icon>
                </div>
                <div class="selector-info">
                  <h3 class="selector-title">✨ Transform with AI (Optional)</h3>
                  <p class="selector-subtitle">Use OpenAI to generate text before sending</p>
                </div>
              </div>

              <div class="service-selection">
                <div v-if="!form.intermediateActionService" class="service-grid">
                  <div
                    class="service-card optional-card"
                    @click="form.intermediateActionService = 'OpenAI'"
                  >
                    <div class="service-card-icon">
                      <img :src="getFallbackIcon('OpenAI')" alt="OpenAI" class="service-icon" />
                    </div>
                    <span class="service-card-name">Add OpenAI</span>
                  </div>
                </div>

                <div v-else class="selected-service-display">
                  <div class="selected-service-card">
                    <div class="service-avatar">
                      <img :src="getFallbackIcon('OpenAI')" alt="OpenAI" class="service-icon" />
                    </div>
                    <div class="service-info">
                      <span class="service-name">OpenAI</span>
                      <span class="service-type">Intermediate Action</span>
                    </div>
                    <button class="change-service-btn" @click="form.intermediateActionService = ''">
                      <v-icon size="16">mdi-close</v-icon>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="form.intermediateActionService" class="connection-arrow">
              <div class="arrow-line"></div>
              <div class="arrow-icon">
                <v-icon size="20" color="#8b5cf6">mdi-arrow-down</v-icon>
              </div>
              <div class="arrow-line"></div>
            </div>

            <div class="service-selector">
              <div class="selector-header">
                <div class="selector-icon">
                  <v-icon size="24" color="#3b82f6">mdi-lightning-bolt</v-icon>
                </div>
                <div class="selector-info">
                  <h3 class="selector-title">Then do this</h3>
                  <p class="selector-subtitle">Choose the action service</p>
                </div>
              </div>

              <div class="service-selection">
                <div v-if="servicesError" class="service-error">
                  <v-icon size="18" color="#ef4444">mdi-alert-circle</v-icon>
                  <span>{{ servicesError }}</span>
                </div>

                <div v-else-if="isLoadingServices" class="service-loading">
                  <div class="spinner"></div>
                  <span>Loading services...</span>
                </div>

                <div v-else>
                  <div v-if="!form.actionService" class="service-grid">
                    <div
                      v-for="item in actionItems.slice(0, 15)"
                      :key="item.value"
                      class="service-card"
                      @click="selectAction(item.value)"
                    >
                      <div class="service-card-icon">
                        <img :src="item.icon" :alt="item.title" class="service-icon" />
                      </div>
                      <span class="service-card-name">{{ item.title }}</span>
                    </div>
                    <div class="service-card more-services" @click="showAllReactionServices = true" v-if="actionItems.length > 15">
                      <div class="service-card-icon">
                        <v-icon size="24" color="#3b82f6">mdi-plus</v-icon>
                      </div>
                      <span class="service-card-name">More...</span>
                    </div>
                  </div>

                  <div v-else class="selected-service-display">
                    <div class="selected-service-card">
                      <div class="service-avatar">
                        <img :src="serviceIcons[form.actionService || ''] || getFallbackIcon(form.actionService || '')" :alt="getServiceName(form.actionService)" class="service-icon" />
                      </div>
                      <div class="service-info">
                        <span class="service-name">{{ getServiceName(form.actionService) }}</span>
                        <span class="service-type">Action Service</span>
                      </div>
                      <button class="change-service-btn" @click="form.actionService = ''">
                        <v-icon size="16">mdi-close</v-icon>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="form.triggerService" class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-cog-outline</v-icon>
            <span class="label-text">Configuration</span>
          </div>

          <div v-if="isCalendarTrigger" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('google-calendar.png')" alt="Google Calendar" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📅 Calendar Event Trigger</h4>
                <p class="config-subtitle">Configure when this area should trigger</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">📅 Event Date & Time</label>
                  <input
                    v-model="form.triggerConfig.eventTime"
                    type="datetime-local"
                    class="modern-input"
                    :min="new Date().toISOString().slice(0, 16)"
                    required
                  />
                  <small class="input-hint">Select the date and time when you want to be reminded</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Event Title (Optional)</label>
                  <input
                    v-model="form.triggerConfig.eventTitle"
                    class="modern-input"
                    placeholder="e.g., Meeting with John, Doctor Appointment"
                  />
                  <small class="input-hint">This will be used in the email subject and body</small>
                </div>

                <div class="input-container" v-if="form.triggerService === 'Google Calendar'">
                  <label class="input-label">🗓️ Calendar ID</label>
                  <input
                    v-model="form.triggerConfig.calendarId"
                    class="modern-input"
                    placeholder="primary"
                  />
                  <small class="input-hint">Use 'primary' for your main calendar</small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'GitHub'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('github.png')" alt="GitHub" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">🐙 GitHub Repository Trigger</h4>
                <p class="config-subtitle">Configure which GitHub events should trigger this area</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">Select Repository</label>
                  <select
                    v-model="form.triggerConfig.repositoryId"
                    class="modern-input"
                    @change="onRepositoryChange"
                    :disabled="isLoadingRepositories"
                    required
                  >
                    <option value="">{{ isLoadingRepositories ? 'Loading repositories...' : 'Choose a repository...' }}</option>
                    <option
                      v-for="repo in repositories"
                      :key="repo.id"
                      :value="repo.id"
                    >
                      {{ repo.full_name }} {{ repo.private ? '(Private)' : '' }}
                    </option>
                  </select>
                  <small class="input-hint">
                    {{ isLoadingRepositories ? 'Loading your GitHub repositories...' : 'Select the repository you want to monitor' }}
                  </small>
                </div>

                <div class="input-container">
                  <label class="input-label">Event Types</label>
                  <div class="checkbox-group">
                    <label class="checkbox-item">
                      <input
                        type="checkbox"
                        v-model="form.triggerConfig.notificationTypes"
                        value="push"
                        @change="onNotificationTypeChange"
                      />
                      <span class="checkbox-label">Push Events</span>
                    </label>
                    <label class="checkbox-item">
                      <input
                        type="checkbox"
                        v-model="form.triggerConfig.notificationTypes"
                        value="pull_request"
                        @change="onNotificationTypeChange"
                      />
                      <span class="checkbox-label">Pull Requests</span>
                    </label>
                    <label class="checkbox-item">
                      <input
                        type="checkbox"
                        v-model="form.triggerConfig.notificationTypes"
                        value="issues"
                        @change="onNotificationTypeChange"
                      />
                      <span class="checkbox-label">Issues</span>
                    </label>
                  </div>
                  <small class="input-hint">Choose which GitHub events should trigger this area</small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'Weather'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('weather.png')" alt="Weather" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">Weather Trigger</h4>
                <p class="config-subtitle">Configure weather conditions that should trigger this area</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">City</label>
                  <select v-model="form.triggerConfig.city" class="modern-select" required>
                    <option value="">Select a city...</option>
                    <option value="Paris">Paris, France</option>
                    <option value="London">London, UK</option>
                    <option value="New York">New York, USA</option>
                    <option value="Tokyo">Tokyo, Japan</option>
                    <option value="Berlin">Berlin, Germany</option>
                    <option value="Madrid">Madrid, Spain</option>
                    <option value="Rome">Rome, Italy</option>
                    <option value="Amsterdam">Amsterdam, Netherlands</option>
                    <option value="Barcelona">Barcelona, Spain</option>
                    <option value="Vienna">Vienna, Austria</option>
                    <option value="Prague">Prague, Czech Republic</option>
                    <option value="Warsaw">Warsaw, Poland</option>
                    <option value="Budapest">Budapest, Hungary</option>
                    <option value="Stockholm">Stockholm, Sweden</option>
                    <option value="Copenhagen">Copenhagen, Denmark</option>
                    <option value="Oslo">Oslo, Norway</option>
                    <option value="Helsinki">Helsinki, Finland</option>
                    <option value="Zurich">Zurich, Switzerland</option>
                    <option value="Brussels">Brussels, Belgium</option>
                    <option value="Dublin">Dublin, Ireland</option>
                    <option value="Lisbon">Lisbon, Portugal</option>
                    <option value="Athens">Athens, Greece</option>
                    <option value="Istanbul">Istanbul, Turkey</option>
                    <option value="Moscow">Moscow, Russia</option>
                    <option value="Sydney">Sydney, Australia</option>
                    <option value="Melbourne">Melbourne, Australia</option>
                    <option value="Toronto">Toronto, Canada</option>
                    <option value="Vancouver">Vancouver, Canada</option>
                    <option value="Montreal">Montreal, Canada</option>
                    <option value="Los Angeles">Los Angeles, USA</option>
                    <option value="Chicago">Chicago, USA</option>
                    <option value="San Francisco">San Francisco, USA</option>
                    <option value="Boston">Boston, USA</option>
                    <option value="Miami">Miami, USA</option>
                    <option value="Seattle">Seattle, USA</option>
                    <option value="Denver">Denver, USA</option>
                    <option value="Las Vegas">Las Vegas, USA</option>
                    <option value="Phoenix">Phoenix, USA</option>
                    <option value="Houston">Houston, USA</option>
                    <option value="Dallas">Dallas, USA</option>
                    <option value="Atlanta">Atlanta, USA</option>
                    <option value="Detroit">Detroit, USA</option>
                    <option value="Philadelphia">Philadelphia, USA</option>
                    <option value="Washington">Washington, USA</option>
                    <option value="Beijing">Beijing, China</option>
                    <option value="Shanghai">Shanghai, China</option>
                    <option value="Hong Kong">Hong Kong, China</option>
                    <option value="Singapore">Singapore</option>
                    <option value="Bangkok">Bangkok, Thailand</option>
                    <option value="Seoul">Seoul, South Korea</option>
                    <option value="Mumbai">Mumbai, India</option>
                    <option value="Delhi">Delhi, India</option>
                    <option value="Bangalore">Bangalore, India</option>
                    <option value="Dubai">Dubai, UAE</option>
                    <option value="Tel Aviv">Tel Aviv, Israel</option>
                    <option value="Cairo">Cairo, Egypt</option>
                    <option value="Johannesburg">Johannesburg, South Africa</option>
                    <option value="Cape Town">Cape Town, South Africa</option>
                    <option value="São Paulo">São Paulo, Brazil</option>
                    <option value="Rio de Janeiro">Rio de Janeiro, Brazil</option>
                    <option value="Buenos Aires">Buenos Aires, Argentina</option>
                    <option value="Santiago">Santiago, Chile</option>
                    <option value="Lima">Lima, Peru</option>
                    <option value="Bogotá">Bogotá, Colombia</option>
                    <option value="Mexico City">Mexico City, Mexico</option>
                  </select>
                  <small class="input-hint">Select the city to monitor weather conditions</small>
                </div>

                <div class="input-container">
                  <label class="input-label">Temperature Threshold</label>
                  <div class="temperature-input-group">
                    <select v-model="form.triggerConfig.operator" class="modern-select operator-select">
                      <option value="greater_than">Greater than</option>
                      <option value="less_than">Less than</option>
                      <option value="equals">Equals</option>
                    </select>
                    <input
                      v-model.number="form.triggerConfig.temperature"
                      type="number"
                      class="modern-input temperature-input"
                      placeholder="25"
                      step="0.1"
                    />
                    <span class="temperature-unit">°C</span>
                  </div>
                  <small class="input-hint">Set the temperature condition to trigger the area</small>
                </div>

                <div class="input-container">
                  <label class="input-label">Weather Condition</label>
                  <select v-model="form.triggerConfig.condition" class="modern-select">
                    <option value="">Any condition</option>
                    <option value="Clear">Clear</option>
                    <option value="Clouds">Cloudy</option>
                    <option value="Rain">Rainy</option>
                    <option value="Snow">Snowy</option>
                    <option value="Thunderstorm">Thunderstorm</option>
                    <option value="Drizzle">Drizzle</option>
                    <option value="Mist">Misty</option>
                    <option value="Fog">Foggy</option>
                  </select>
                  <small class="input-hint">Optional: Trigger only for specific weather conditions</small>
                </div>
              </div>

              <div class="test-section">
                <button
                  type="button"
                  class="test-weather-btn"
                  @click="testWeatherTrigger"
                  :disabled="!form.triggerConfig.city || isTestingWeather"
                >
                  <v-icon size="18">mdi-flash</v-icon>
                  {{ isTestingWeather ? 'Testing...' : 'Test Weather Trigger' }}
                </button>
                <div v-if="weatherTestResult" class="test-result">
                  <div :class="['test-status', weatherTestResult.triggered ? 'success' : 'info']">
                    <v-icon size="16">{{ weatherTestResult.triggered ? 'mdi-check-circle' : 'mdi-information' }}</v-icon>
                    {{ weatherTestResult.message }}
                  </div>
                  <div v-if="weatherTestResult.data" class="weather-data">
                    <small>Current: {{ weatherTestResult.data.temperature }}°C, {{ weatherTestResult.data.condition }}</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'Google Drive'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('google-drive.png')" alt="Google Drive" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📁 Google Drive Trigger</h4>
                <p class="config-subtitle">Trigger when a new file appears in a folder</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">🗂️ Folder ID</label>
                  <input
                    v-model="form.triggerConfig.folderId"
                    class="modern-input"
                    placeholder="e.g., 1A2B3C4D..."
                    required
                  />
                  <small class="input-hint">Copy the ID from the Drive folder URL (between /folders/ and the end)</small>
                </div>

                <div class="input-container">
                  <label class="input-label">🔖 Email Subject (Gmail action)</label>
                  <input
                    v-if="form.actionService === 'Gmail'"
                    v-model="form.actionConfig.subject"
                    class="modern-input"
                    placeholder="New Drive file: {{fileName}}"
                  />
                  <small v-if="form.actionService === 'Gmail'" class="input-hint" v-pre>You can use {{fileName}}, {{mimeType}}, {{webViewLink}}</small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'Google Sheets'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('google-sheets.png')" alt="Google Sheets" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📊 Google Sheets Trigger</h4>
                <p class="config-subtitle">Surveillez une plage de votre feuille et déclenchez des actions sur chaque modification</p>
              </div>
            </div>

            <div class="config-content">
              <div class="sheets-config-grid">
                <div class="input-container">
                  <label class="input-label">🆔 Spreadsheet ID</label>
                  <input
                    v-model="form.triggerConfig.spreadsheetId"
                    class="modern-input"
                    placeholder="1A2B3C4D..."
                    required
                  />
                  <small class="input-hint">Copiez l'identifiant présent dans l'URL de votre feuille (entre /d/ et /edit)</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📄 Nom de la feuille (optionnel)</label>
                  <input
                    v-model="form.triggerConfig.sheetName"
                    class="modern-input"
                    placeholder="Feuille1"
                  />
                  <small class="input-hint">Utilisé uniquement pour vos logs et messages</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📍 Plage A1</label>
                  <input
                    v-model="form.triggerConfig.range"
                    class="modern-input"
                    placeholder="Feuille1!A1:D"
                    required
                  />
                  <small class="input-hint">Définissez la plage à surveiller (format A1). Limitez-la pour de meilleures performances.</small>
                </div>

                <label class="sheets-checkbox">
                  <input
                    v-model="form.triggerConfig.hasHeader"
                    type="checkbox"
                  />
                  <span>La première ligne contient des en-têtes</span>
                </label>
              </div>

              <div class="sheets-test-actions">
                <button
                  type="button"
                  class="test-google-sheets-btn"
                  @click="testGoogleSheets"
                  :disabled="!canTestGoogleSheets || isTestingGoogleSheets"
                >
                  <v-icon size="18">{{ isTestingGoogleSheets ? 'mdi-loading' : 'mdi-table-arrow-down' }}</v-icon>
                  {{ isTestingGoogleSheets ? 'Test en cours...' : 'Tester la connexion' }}
                </button>
                <div v-if="sheetsTestError" class="error-message">
                  ❌ {{ sheetsTestError }}
                </div>
                <div v-else-if="sheetsTestResult" class="sheets-test-result">
                  <div class="sheets-test-summary">
                    <v-icon size="16" color="#22c55e">mdi-check-circle</v-icon>
                    <span>{{ sheetsTestResult.rowCount }} lignes récupérées</span>
                  </div>
                  <div v-if="sheetsTestResult.previewRows.length" class="sheets-test-preview">
                    <table>
                      <tbody>
                        <tr v-for="(row, rowIndex) in sheetsTestResult.previewRows" :key="rowIndex">
                          <td v-for="(cell, cellIndex) in row" :key="cellIndex">
                            {{ cell }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <small class="input-hint">Aperçu limité aux 5 premières lignes</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'Timer'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('google-calendar.png')" alt="Timer" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">⏰ Timer Trigger</h4>
                <p class="config-subtitle">Configure the interval for automatic execution</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">⏱️ Time Interval</label>
                  <select v-model="form.triggerConfig.interval" class="modern-select" required>
                    <option value="">Select an interval...</option>
                    <option value="30s">Every 30 seconds (testing)</option>
                    <option value="1m">Every 1 minute</option>
                    <option value="5m">Every 5 minutes</option>
                    <option value="15m">Every 15 minutes</option>
                    <option value="30m">Every 30 minutes</option>
                    <option value="1h">Every 1 hour</option>
                    <option value="2h">Every 2 hours</option>
                    <option value="6h">Every 6 hours</option>
                    <option value="12h">Every 12 hours</option>
                    <option value="24h">Every 24 hours (daily)</option>
                    <option value="168h">Every 7 days (weekly)</option>
                  </select>
                  <small class="input-hint">Choose how often this area should trigger automatically</small>
                </div>

                <div v-if="form.triggerConfig.interval" class="input-container">
                  <div class="info-box">
                    <v-icon size="18" color="#3b82f6">mdi-information</v-icon>
                    <span>Your area will run automatically every {{ form.triggerConfig.interval }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.intermediateActionService === 'OpenAI'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getFallbackIcon('OpenAI')" alt="OpenAI" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">🤖 OpenAI Configuration</h4>
                <p class="config-subtitle">Configure how OpenAI should generate text from your trigger data</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">💬 System Prompt (Optional)</label>
                  <textarea
                    v-model="form.intermediateActionConfig.systemPrompt"
                    class="modern-textarea"
                    placeholder="You are a helpful assistant..."
                    rows="3"
                  ></textarea>
                  <small class="input-hint">Define the AI's role or behavior. This helps guide the response style.</small>
                </div>

                <div class="input-container">
                  <label class="input-label">✍️ Prompt</label>
                  <textarea
                    v-model="form.intermediateActionConfig.prompt"
                    class="modern-textarea"
                    placeholder="Generate a summary of the following data: {{changeType}} in {{sheetName}} at row {{rowNumber}}"
                    rows="5"
                    required
                  ></textarea>
                  <small class="input-hint">
                    Use template variables from your trigger:
                    <span v-if="form.triggerService === 'Telegram'">&#123;&#123;messageText&#125;&#125;, &#123;&#123;firstName&#125;&#125;, &#123;&#123;username&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'Timer'">&#123;&#123;triggerTime&#125;&#125;, &#123;&#123;interval&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'Google Sheets'">&#123;&#123;changeType&#125;&#125;, &#123;&#123;sheetName&#125;&#125;, &#123;&#123;rowNumber&#125;&#125;, &#123;&#123;rowData&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'GitHub'">&#123;&#123;repository&#125;&#125;, &#123;&#123;branch&#125;&#125;</span>
                    <span v-else>&#123;&#123;areaName&#125;&#125;, &#123;&#123;eventTime&#125;&#125;</span>
                  </small>
                  <div class="info-box" style="margin-top: 0.5rem;">
                    <v-icon size="16" color="#8b5cf6">mdi-lightbulb-on</v-icon>
                    <span>The generated text will be available as <code>&#123;&#123;openaiGeneratedText&#125;&#125;</code> in your action</span>
                  </div>
                </div>

                <div class="input-group" style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                  <div class="input-container">
                    <label class="input-label">🌡️ Temperature</label>
                    <input
                      v-model.number="form.intermediateActionConfig.temperature"
                      type="number"
                      min="0"
                      max="2"
                      step="0.1"
                      class="modern-input"
                    />
                    <small class="input-hint">Creativity level (0 = focused, 2 = creative). Default: 0.7</small>
                  </div>

                  <div class="input-container">
                    <label class="input-label">📏 Max Tokens</label>
                    <input
                      v-model.number="form.intermediateActionConfig.maxTokens"
                      type="number"
                      min="1"
                      max="4000"
                      step="50"
                      class="modern-input"
                    />
                    <small class="input-hint">Maximum length of generated text. Default: 500</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'Telegram'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('telegram.png')" alt="Telegram" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">💬 Telegram Trigger</h4>
                <p class="config-subtitle">Configure which Telegram messages should trigger this area</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">📱 Chat ID</label>
                  <input
                    v-model="form.triggerConfig.chatId"
                    class="modern-input"
                    placeholder="123456789"
                    required
                  />
                  <small class="input-hint">Your Telegram chat ID. Send a message to your bot and visit the Telegram API to get it.</small>
                </div>

                <div class="input-container">
                  <label class="input-label">🎯 Trigger Type</label>
                  <select v-model="form.triggerConfig.triggerType" class="modern-select" required>
                    <option value="">Select trigger type...</option>
                    <option value="message_received">Any Message Received</option>
                    <option value="keyword_match">Keyword Match</option>
                    <option value="command_received">Command Received</option>
                  </select>
                  <small class="input-hint">Choose when this area should trigger</small>
                </div>

                <div v-if="form.triggerConfig.triggerType === 'keyword_match'" class="input-container">
                  <label class="input-label">🔑 Keyword</label>
                  <input
                    v-model="form.triggerConfig.keyword"
                    class="modern-input"
                    placeholder="urgent"
                    required
                  />
                  <small class="input-hint">Messages containing this keyword will trigger the area</small>
                </div>

                <div v-if="form.triggerConfig.triggerType === 'command_received'" class="input-container">
                  <label class="input-label">⚡ Command</label>
                  <input
                    v-model="form.triggerConfig.command"
                    class="modern-input"
                    placeholder="/start"
                    required
                  />
                  <small class="input-hint">This specific command will trigger the area (e.g., /start, /help)</small>
                </div>

                <div v-if="form.triggerConfig.chatId" class="input-container">
                  <div class="info-box">
                    <v-icon size="18" color="#3b82f6">mdi-information</v-icon>
                    <span>
                      <span v-if="form.triggerConfig.triggerType === 'message_received'">
                        Any message sent to this chat will trigger the area
                      </span>
                      <span v-else-if="form.triggerConfig.triggerType === 'keyword_match'">
                        Only messages containing "{{ form.triggerConfig.keyword }}" will trigger
                      </span>
                      <span v-else-if="form.triggerConfig.triggerType === 'command_received'">
                        Only the command "{{ form.triggerConfig.command }}" will trigger
                      </span>
                    </span>
                  </div>
                </div>

                <div class="input-container">
                  <div class="info-box warning">
                    <v-icon size="18" color="#f59e0b">mdi-alert</v-icon>
                    <span>Make sure you have set up the Telegram webhook pointing to your backend!</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.actionService === 'Gmail'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('gmail.png')" alt="Gmail" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📧 Gmail Action</h4>
                <p class="config-subtitle">Configure the email to be sent</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">📧 Send Email To</label>
                  <input
                    v-model="form.actionConfig.toEmail"
                    type="email"
                    class="modern-input"
                    placeholder="your-email@gmail.com"
                    required
                  />
                  <small class="input-hint">Enter the email address where you want to receive the reminder</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Email Subject</label>
                  <input
                    v-model="form.actionConfig.subject"
                    class="modern-input"
                    placeholder="Reminder: {{eventTitle}}"
                    required
                  />
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125; to include the event name<span v-if="form.intermediateActionService === 'OpenAI'">, or &#123;&#123;openaiGeneratedText&#125;&#125; for AI-generated content</span></small>
                </div>

                <div class="input-container">
                  <label class="input-label">📄 Email Body</label>
                  <textarea
                    v-model="form.actionConfig.body"
                    class="modern-textarea"
                    placeholder="Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}"
                    rows="4"
                    required
                  ></textarea>
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125;, &#123;&#123;eventTime&#125;&#125;, and &#123;&#123;areaName&#125;&#125; as placeholders<span v-if="form.intermediateActionService === 'OpenAI'">, or &#123;&#123;openaiGeneratedText&#125;&#125; for AI-generated content</span></small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.actionService === 'Discord'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('discord.png')" alt="Discord" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">💬 Discord Action</h4>
                <p class="config-subtitle">Configure the Discord message to be sent</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">🔗 Discord Webhook URL</label>
                  <input
                    v-model="form.actionConfig.webhookUrl"
                    class="modern-input"
                    placeholder="https://discord.com/api/webhooks/..."
                    required
                  />
                  <small class="input-hint">Your Discord webhook URL (Server Settings → Integrations → Webhooks)</small>
                </div>

                <div class="input-container">
                  <label class="input-label">💬 Message</label>
                  <textarea
                    v-model="form.actionConfig.message"
                    class="modern-textarea"
                    placeholder="Message to send to Discord"
                    rows="4"
                    required
                  ></textarea>
                  <small class="input-hint">
                    Use template variables:
                    <span v-if="form.triggerService === 'Telegram'">&#123;&#123;messageText&#125;&#125;, &#123;&#123;firstName&#125;&#125;, &#123;&#123;username&#125;&#125;, &#123;&#123;chatId&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'Timer'">&#123;&#123;triggerTime&#125;&#125;, &#123;&#123;interval&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'Google Sheets'">&#123;&#123;changeType&#125;&#125;, &#123;&#123;sheetName&#125;&#125;, &#123;&#123;rowNumber&#125;&#125;, &#123;&#123;rowData&#125;&#125;</span>
                    <span v-else-if="form.triggerService === 'Spotify'">&#123;&#123;trackName&#125;&#125;, &#123;&#123;artistNames&#125;&#125;, &#123;&#123;albumName&#125;&#125;, &#123;&#123;trackUrl&#125;&#125;</span>
                    <span v-else>&#123;&#123;areaName&#125;&#125;, &#123;&#123;eventTime&#125;&#125;</span>
                    <span v-if="form.intermediateActionService === 'OpenAI'">, &#123;&#123;openaiGeneratedText&#125;&#125;</span>
                  </small>
                </div>

                <div v-if="form.actionConfig.message" class="input-container">
                  <div class="info-box">
                    <v-icon size="18" color="#3b82f6">mdi-information</v-icon>
                    <span>This message will be posted to your Discord channel via webhook</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.actionService === 'Telegram'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('telegram.png')" alt="Telegram" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📱 Telegram Action</h4>
                <p class="config-subtitle">Configure the Telegram message to be sent</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">💬 Chat ID</label>
                  <input
                    v-model="form.actionConfig.chatId"
                    type="text"
                    class="modern-input"
                    placeholder="8481009224"
                    required
                  />
                  <small class="input-hint">
                    Your Telegram chat ID.
                    <a href="https://t.me/userinfobot" target="_blank" class="helper-link">
                      Get it from @userinfobot
                    </a>
                    (send /start to the bot)
                  </small>
                  <div class="help-box">
                    <v-icon size="16" color="#22c55e">mdi-help-circle</v-icon>
                    <div class="help-content">
                      <strong>How to find your Chat ID:</strong>
                      <ol>
                        <li>Open Telegram and search for <strong>@userinfobot</strong></li>
                        <li>Click on it and send <code>/start</code></li>
                        <li>The bot will reply with your ID (ex: 987654321)</li>
                        <li>Copy and paste it here</li>
                      </ol>
                    </div>
                  </div>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Message</label>
                  <textarea
                    v-model="form.actionConfig.message"
                    class="modern-textarea"
                    placeholder="🤖 Notification from {{areaName}}&#10;⏰ Triggered at {{triggerTime}}"
                    rows="5"
                    required
                  ></textarea>
                  <small class="input-hint">
                    Use template variables: &#123;&#123;areaName&#125;&#125;, &#123;&#123;triggerTime&#125;&#125;, &#123;&#123;interval&#125;&#125;, etc.
                    <span v-if="form.intermediateActionService === 'OpenAI'"> Use &#123;&#123;openaiGeneratedText&#125;&#125; to include the AI-generated text.</span>
                    Telegram supports Markdown formatting (*bold*, _italic_)
                  </small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="isCalendarTrigger && form.actionService === 'Gmail'" class="preview-section">
            <div class="preview-header">
              <v-icon class="preview-icon" size="20">mdi-eye-outline</v-icon>
              <span class="preview-title">Email Preview</span>
            </div>
            <div class="preview-content">
              <div class="email-preview">
                <div class="email-header">
                  <strong>To:</strong> {{ form.actionConfig.toEmail || 'your-email@gmail.com' }}
                </div>
                <div class="email-header">
                  <strong>Subject:</strong> {{ form.actionConfig.subject || 'Reminder: Event Title' }}
                </div>
                <div class="email-body">
                  {{ form.actionConfig.body || 'Hello! This is a reminder about your upcoming event: Event Title at Event Time' }}
                </div>
              </div>
            </div>
          </div>

          <div v-if="form.triggerService === 'GitHub' && form.actionService === 'Gmail'" class="preview-section">
            <div class="preview-header">
              <v-icon class="preview-icon" size="20">mdi-eye-outline</v-icon>
              <span class="preview-title">GitHub → Gmail Preview</span>
            </div>
            <div class="preview-content">
              <div class="github-preview">
                <div class="preview-item">
                  <strong>Repository:</strong> {{ getSelectedRepositoryName() || 'Select a repository' }}
                </div>
                <div class="preview-item">
                  <strong>Events:</strong> {{ form.triggerConfig.notificationTypes?.join(', ') || 'Select event types' }}
                </div>
                <div class="preview-item">
                  <strong>Email To:</strong> {{ form.actionConfig.toEmail || 'your-email@gmail.com' }}
                </div>
                <div class="preview-item">
                  <strong>Subject:</strong> {{ form.actionConfig.subject || 'GitHub Activity Notification' }}
                </div>
                <div class="preview-item">
                  <strong>Body:</strong> {{ form.actionConfig.body || 'New activity detected in your repository' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="error" class="error-message">
        <v-icon size="20" color="#ef4444">mdi-alert-circle</v-icon>
        <span>{{ error }}</span>
      </div>

      <div class="card-actions">
        <button class="action-btn secondary" @click="$emit('close')">
          <v-icon size="18">mdi-close</v-icon>
          Cancel
        </button>
        <button class="action-btn primary" @click="createArea" :disabled="!isFormValid || isLoading">
          <v-icon size="18">mdi-check</v-icon>
          {{ isLoading ? 'Creating...' : 'Create Area' }}
        </button>

        <button
          v-if="isCalendarTrigger && form.actionService === 'Gmail'"
          class="action-btn test-email-btn"
          @click="sendTestEmail"
          :disabled="!canSendTestEmail || isSendingTest"
        >
          <v-icon size="18">mdi-email-send</v-icon>
          {{ isSendingTest ? 'Sending...' : 'Send Test Email' }}
        </button>

        <button
          v-if="form.triggerService === 'GitHub' && form.actionService === 'Gmail'"
          class="action-btn test-email-btn"
          @click="sendTestEmail"
          :disabled="!canSendTestEmail || isSendingTest"
        >
          <v-icon size="18">mdi-email-send</v-icon>
          {{ isSendingTest ? 'Sending...' : 'Send Test Email' }}
        </button>

        <div v-if="isCalendarTrigger && form.actionService === 'Gmail'" class="debug-info">
          <small style="color: #666; font-size: 0.75rem;">
            Debug: {{ isFormValid ? 'Ready to create' : 'Missing: ' + getMissingFields() }}
          </small>
        </div>

        <div v-if="form.triggerService === 'GitHub' && form.actionService === 'Gmail'" class="debug-info">
          <small style="color: #666; font-size: 0.75rem;">
            Debug: {{ isFormValid ? 'Ready to create' : 'Missing: ' + getMissingFields() }}
          </small>
        </div>
      </div>
    </div>

    <AreaGuideModal :is-open="showGuide" @close="showGuide = false" />
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { API_BASE_URL } from '../../config/api'
import { areaService, type GoogleSheetsTestResponse } from '../../services/area'
import { githubService, type GitHubRepository } from '../../services/github'
import { useAuth } from '@/composables/useAuth'
import AreaGuideModal from '../AreaGuideModal.vue'

type ServiceInfo = {
  name: string
  actions: { name: string; description: string }[]
  reactions: { name: string; description: string }[]
}

interface AreaTemplate {
  id: string
  title: string
  subtitle: string
  description: string
  icon: string
  gradientClass: string
  triggerService: string
  actionService: string
  isActive: boolean
}

const props = defineProps<{
  template?: AreaTemplate | null
}>()

const { currentUser } = useAuth()
const ICONS_DIR = 'app-icons'

const getIconUrl = (file: string) =>
  new URL(`../../assets/${ICONS_DIR}/${file}`, import.meta.url).href

const priorityServices = [
  'Google Calendar',
  'Date Timer',
  'Gmail',
  'Weather',
  'Discord',
  'Google Drive',
  'Google Sheets',
  'GitHub',
  'Timer',
  'Telegram',
  'Spotify',
  'Twitter'
]

const intermediateServices = [
  'OpenAI'
]

const CALENDAR_SERVICES = ['Google Calendar', 'Date Timer']

const isCalendarService = (service?: string | null) =>
  CALENDAR_SERVICES.includes(service ?? '')

const services = ref<ServiceInfo[]>([])
const serviceIcons = ref<Record<string, string>>({})
const isLoadingServices = ref(false)
const servicesError = ref<string | null>(null)

const fetchServices = async () => {
  isLoadingServices.value = true
  servicesError.value = null

  try {
    const response = await fetch(`${API_BASE_URL}/about.json`)
    if (!response.ok) {
      throw new Error(`Failed to fetch services: ${response.status}`)
    }

    const data = await response.json()
    const fetchedServices = (data?.server?.services || []) as ServiceInfo[]

    services.value = fetchedServices

    const icons: Record<string, string> = {}
    const defaultIcons: Record<string, string> = {
      Gmail: 'gmail.png',
      Slack: 'slack.png',
      GitHub: 'github.png',
      Weather: 'weather.png',
      'Google Calendar': 'google-calendar.png',
      Discord: 'discord.png',
      'Google Sheets': 'google-sheets.png',
      'Google Drive': 'google-drive.png',
      Timer: 'google-calendar.png',
      'Date Timer': 'google-calendar.png',
      Telegram: 'telegram.png',
      OpenAI: 'openai.png'
    }

    fetchedServices.forEach(service => {
      const key = service.name
      const normalized = key.toLowerCase()
      const matchedDefault = Object.keys(defaultIcons).find(name => name.toLowerCase() === normalized)
      if (matchedDefault) {
        icons[key] = getIconUrl(defaultIcons[matchedDefault])
      } else {
        const sanitized = key.toLowerCase().replace(/\s+/g, '-')
        icons[key] = getIconUrl(`${sanitized}.png`)
      }
    })

    serviceIcons.value = icons
  } catch (err) {
    servicesError.value = err instanceof Error ? err.message : 'Failed to load services'
  } finally {
    isLoadingServices.value = false
  }
}

onMounted(() => {
  fetchServices()
})

const getFallbackIcon = (serviceName: string) => {
  if (!serviceName) {
    return getIconUrl('gmail.png')
  }

  const normalized = serviceName.toLowerCase()
  if (serviceIcons.value[serviceName]) {
    return serviceIcons.value[serviceName]
  }

  const defaultIcons: Record<string, string> = {
    gmail: 'gmail.png',
    slack: 'slack.png',
    github: 'github.png',
    weather: 'weather.png',
    'google calendar': 'google-calendar.png',
    'date timer': 'google-calendar.png',
    discord: 'discord.png',
    'google sheets': 'google-sheets.png',
    'google drive': 'google-drive.png',
    timer: 'google-calendar.png',
    telegram: 'telegram.png',
    openai: 'openai.png',
    spotify: 'spotify.png'
  }

  const matchedDefault = Object.keys(defaultIcons).find(name => name === normalized)
  if (matchedDefault) {
    return getIconUrl(defaultIcons[matchedDefault])
  }

  const sanitized = normalized.replace(/\s+/g, '-')
  return getIconUrl(`${sanitized}.png`)
}

const appItems = computed(() => {
  if (services.value.length === 0) {
    return []
  }

  const sortedServices = [...services.value].sort((a, b) => {
    const aPriority = priorityServices.indexOf(a.name)
    const bPriority = priorityServices.indexOf(b.name)

    if (aPriority === -1 && bPriority === -1) {
      return a.name.localeCompare(b.name)
    }
    if (aPriority === -1) return 1
    if (bPriority === -1) return -1
    return aPriority - bPriority
  })

  return sortedServices.map(service => ({
    title: service.name,
    value: service.name,
    icon: serviceIcons.value[service.name] || getFallbackIcon(service.name)
  }))
})

const triggerItems = computed(() => {
  return appItems.value.filter(item => item.value !== 'OpenAI')
})

const actionItems = computed(() => {
  return appItems.value.filter(item => item.value !== 'OpenAI')
})

const form = reactive({
  areaName: '',
  description: '',
  triggerService: '' as string | null,
  actionService: '' as string | null,
  intermediateActionService: '' as string | null,
  triggerConfig: {} as any,
  actionConfig: {} as any,
  intermediateActionConfig: {
    prompt: '',
    systemPrompt: '',
    temperature: 0.7,
    maxTokens: 500
  } as any,
})

const isCalendarTrigger = computed(() => isCalendarService(form.triggerService))

const repositories = ref<GitHubRepository[]>([])
const isLoadingRepositories = ref(false)
const isTestingGoogleSheets = ref(false)
const sheetsTestError = ref<string | null>(null)
const sheetsTestResult = ref<GoogleSheetsTestResponse | null>(null)
const showGuide = ref(false)
watch(() => props.template, (newTemplate) => {
  if (newTemplate) {
    form.areaName = newTemplate.title
    form.description = newTemplate.description
    form.triggerService = newTemplate.triggerService
    form.actionService = newTemplate.actionService

    if (isCalendarService(newTemplate.triggerService) && newTemplate.actionService === 'Gmail') {
      form.triggerConfig = {
        eventTime: '',
        eventTitle: '',
        calendarId: 'primary'
      }
      form.actionConfig = {
        toEmail: '',
        subject: 'Reminder: {{eventTitle}}',
        body: 'Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}'
      }
    } else if (newTemplate.triggerService === 'GitHub') {
      form.triggerConfig = {
        repository: '',
        branch: '',
        events: [],
        webhookSecret: ''
      }
    } else if (newTemplate.triggerService === 'Weather') {
      form.triggerConfig = {
        city: '',
        temperature: 30,
        condition: ''
      }
    } else if (newTemplate.triggerService === 'Google Sheets') {
      form.triggerConfig = {
        spreadsheetId: '',
        sheetName: '',
        range: 'Sheet1!A1:D',
        hasHeader: true
      }
    } else if (newTemplate.triggerService === 'Telegram') {
      form.triggerConfig = {
        chatId: '',
        triggerType: 'message_received',
        keyword: '',
        command: ''
      }
    }

    if (newTemplate.triggerService === 'GitHub' && newTemplate.actionService === 'Gmail') {
      form.triggerConfig = {
        repositoryId: '',
        notificationTypes: ['push']
      }
      form.actionConfig = {
        toEmail: '',
        subject: 'GitHub Activity: {{repository_name}}',
        body: 'New {{event_type}} activity detected in repository {{repository_name}}.\n\nDetails: {{event_details}}\n\nArea: {{areaName}}'
      }
    } else if (newTemplate.actionService === 'Spotify') {
      const hasHeader = typeof form.triggerConfig?.hasHeader === 'boolean' ? form.triggerConfig.hasHeader : true
      form.actionConfig = {
        playlistId: '',
        spreadsheetId: form.triggerConfig?.spreadsheetId || '',
        range: form.triggerConfig?.range || 'Sheet1!A2:C',
        urlColumn: 'SpotifyLink',
        hasHeader
      }
    }
  }
}, { immediate: true })

const isFormValid = computed(() => {
  const hasBasicInfo = form.areaName.trim() !== '' &&
                      form.triggerService !== '' &&
                      form.actionService !== ''

  if (isCalendarService(form.triggerService) && form.actionService === 'Gmail') {
    return hasBasicInfo &&
           form.triggerConfig.eventTime &&
           form.actionConfig.toEmail &&
           form.actionConfig.subject
  }

  if (form.triggerService === 'GitHub' && form.actionService === 'Gmail') {
    return hasBasicInfo &&
           form.triggerConfig.repositoryId &&
           form.triggerConfig.notificationTypes?.length > 0 &&
           form.actionConfig.toEmail &&
           form.actionConfig.subject
  }

  if (form.triggerService === 'Google Drive') {
    return hasBasicInfo && form.triggerConfig.folderId
  }

  if (form.triggerService === 'Weather') {
    return hasBasicInfo &&
           form.triggerConfig.city &&
           form.triggerConfig.temperature !== undefined
  }

  if (form.triggerService === 'Google Sheets') {
    return hasBasicInfo &&
           form.triggerConfig.spreadsheetId &&
           form.triggerConfig.range
  }

  if (form.triggerService === 'Timer') {
    return hasBasicInfo &&
           form.triggerConfig.interval
  }

  if (form.triggerService === 'Telegram') {
    const hasBasicTelegramConfig = hasBasicInfo &&
           form.triggerConfig.chatId &&
           form.triggerConfig.triggerType

    if (form.triggerConfig.triggerType === 'keyword_match') {
      return hasBasicTelegramConfig && form.triggerConfig.keyword
    }
    if (form.triggerConfig.triggerType === 'command_received') {
      return hasBasicTelegramConfig && form.triggerConfig.command
    }
    return hasBasicTelegramConfig
  }

  if (form.actionService === 'Discord') {
    return hasBasicInfo &&
           form.actionConfig.webhookUrl &&
           form.actionConfig.message
  }

  if (form.actionService === 'Spotify') {
    // Allow creation without config - will be configured later in ConfigureAreaPage
    return hasBasicInfo
  }

  if (form.actionService === 'Telegram') {
    return hasBasicInfo &&
           form.actionConfig.chatId &&
           form.actionConfig.message
  }

  if (form.intermediateActionService === 'OpenAI') {
    return hasBasicInfo &&
           form.intermediateActionConfig.prompt?.trim() !== ''
  }

  return hasBasicInfo
})

const showAllTriggerServices = ref(false)
const showAllReactionServices = ref(false)

const canTestGoogleSheets = computed(() => {
  if (form.triggerService !== 'Google Sheets') {
    return false
  }

  const spreadsheetId = (form.triggerConfig?.spreadsheetId || '').toString().trim()
  const range = (form.triggerConfig?.range || '').toString().trim()

  return !!spreadsheetId && !!range
})

watch(
  () => form.triggerService === 'Google Sheets'
    ? [form.triggerConfig?.spreadsheetId, form.triggerConfig?.range, form.triggerConfig?.sheetName, form.triggerConfig?.hasHeader]
    : null,
  () => {
    if (form.triggerService === 'Google Sheets') {
      sheetsTestResult.value = null
      sheetsTestError.value = null
    }
  }
)

watch(
  () => form.triggerConfig?.spreadsheetId,
  (newValue, oldValue) => {
    if (form.actionService === 'Spotify') {
      const current = (form.actionConfig?.spreadsheetId || '').toString().trim()
      if (!current || current === (oldValue || '').toString().trim()) {
        form.actionConfig.spreadsheetId = (newValue || '').toString().trim()
      }
    }
  }
)

watch(
  () => form.triggerConfig?.range,
  (newValue, oldValue) => {
    if (form.actionService === 'Spotify') {
      const current = (form.actionConfig?.range || '').toString().trim()
      if (!current || current === (oldValue || '').toString().trim()) {
        form.actionConfig.range = (newValue || '').toString().trim()
      }
    }
  }
)

watch(
  () => form.triggerConfig?.hasHeader,
  (newValue) => {
    if (form.actionService === 'Spotify' && typeof newValue === 'boolean') {
      form.actionConfig.hasHeader = newValue
    }
  }
)

const selectTrigger = (serviceId: string) => {
  form.triggerService = serviceId
  showAllTriggerServices.value = false
  sheetsTestResult.value = null
  sheetsTestError.value = null

  if (isCalendarService(serviceId)) {
    form.triggerConfig = {
      eventTime: '',
      eventTitle: '',
      calendarId: 'primary'
    }
  } else if (serviceId === 'GitHub') {
    form.triggerConfig = {
      repository: '',
      branch: '',
      events: [],
      webhookSecret: ''
    }
  } else if (serviceId === 'Weather') {
    form.triggerConfig = {
      city: '',
      temperature: 30,
      condition: ''
    }
  } else if (serviceId === 'Google Sheets') {
    form.triggerConfig = {
      spreadsheetId: '',
      sheetName: '',
      range: 'Sheet1!A1:D',
      hasHeader: true
    }
  } else if (serviceId === 'Google Drive') {
    form.triggerConfig = {
      folderId: '',
      knownFileIds: {},
      lastChecked: null
    }
  } else if (serviceId === 'Timer') {
    form.triggerConfig = {
      interval: '5m'
    }
  } else if (serviceId === 'Telegram') {
    form.triggerConfig = {
      chatId: '',
      triggerType: 'message_received',
      keyword: '',
      command: ''
    }
  } else {
    form.triggerConfig = {}
  }

  if (serviceId === 'GitHub') {
    loadRepositories()
  }
}

const selectAction = (serviceId: string) => {
  form.actionService = serviceId
  showAllReactionServices.value = false

  if (serviceId === 'Gmail') {
    if (form.triggerService === 'GitHub') {
      form.actionConfig = {
        toEmail: '',
        subject: 'GitHub Activity: {{repository_name}}',
        body: 'New {{event_type}} activity detected in repository {{repository_name}}.\n\nDetails: {{event_details}}\n\nArea: {{areaName}}'
      }
    } else {
      form.actionConfig = {
        toEmail: '',
        subject: 'Reminder: {{eventTitle}}',
        body: 'Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}'
      }
    }
  } else if (serviceId === 'Discord') {
    const defaultMessage = form.triggerService === 'Telegram'
      ? '📱 **Telegram Message**\n👤 From: {{firstName}} (@{{username}})\n💬 Message: {{messageText}}\n📱 Chat: {{chatId}}'
      : form.triggerService === 'Google Sheets'
      ? '📊 Google Sheets update ({{changeType}}) in {{sheetName}} row {{rowNumber}}: {{rowData}}'
      : form.triggerService === 'Spotify'
      ? '🎧 Now playing: {{trackName}} — {{artistNames}}\n🔗 {{trackUrl}}'
      : form.triggerService === 'Timer'
      ? '⏰ Timer triggered for {{areaName}}\n📅 Time: {{triggerTime}}\n⏱️ Interval: {{interval}}'
      : 'Automation triggered for {{areaName}}'

    form.actionConfig = {
      webhookUrl: '',
      message: defaultMessage
    }
  } else if (serviceId === 'Spotify') {
    const fallbackRange = form.triggerConfig?.range || 'Sheet1!A2:C'
    const fallbackSheetId = form.triggerConfig?.spreadsheetId || ''
    const hasHeader = typeof form.triggerConfig?.hasHeader === 'boolean' ? form.triggerConfig.hasHeader : true

    form.actionConfig = {
      playlistId: '',
      spreadsheetId: fallbackSheetId,
      range: fallbackRange,
      urlColumn: 'SpotifyLink',
      hasHeader
    }
  } else if (serviceId === 'Telegram') {
    const defaultMessage = form.triggerService === 'Timer'
      ? '⏰ Timer triggered for {{areaName}}\n📅 Time: {{triggerTime}}\n⏱️ Interval: {{interval}}'
      : form.triggerService === 'Google Sheets'
      ? '📊 Google Sheets update ({{changeType}}) in {{sheetName}} row {{rowNumber}}: {{rowData}}'
      : form.triggerService === 'Telegram'
      ? '💬 Telegram message received!\n👤 From: {{firstName}} (@{{username}})\n📝 Message: {{messageText}}\n📱 Chat: {{chatId}}'
      : form.triggerService === 'Spotify'
      ? '🎧 Now playing on Spotify: {{trackName}} — {{artistNames}}'
      : '🤖 Notification from {{areaName}}\n⏰ Triggered at {{triggerTime}}'

    form.actionConfig = {
      chatId: '',
      message: defaultMessage
    }
  } else {
    form.actionConfig = {}
  }
}

const getServiceName = (serviceId: string) => {
  const service = appItems.value.find(item => item.value === serviceId)
  return service?.title || ''
}

const getMissingFields = () => {
  const missing = []
  if (!form.areaName.trim()) missing.push('Area Name')

  if (isCalendarService(form.triggerService)) {
    if (!form.triggerConfig.eventTime) missing.push('Event Time')
  }

  if (form.triggerService === 'GitHub') {
    if (!form.triggerConfig.repositoryId) missing.push('Repository')
    if (!form.triggerConfig.notificationTypes?.length) missing.push('Event Types')
  }

  if (form.triggerService === 'Google Sheets') {
    if (!form.triggerConfig.spreadsheetId) missing.push('Spreadsheet ID')
    if (!form.triggerConfig.range) missing.push('Range')
  }

  if (form.triggerService === 'Telegram') {
    if (!form.triggerConfig.chatId) missing.push('Chat ID')
    if (!form.triggerConfig.triggerType) missing.push('Trigger Type')
    if (form.triggerConfig.triggerType === 'keyword_match' && !form.triggerConfig.keyword) {
      missing.push('Keyword')
    }
    if (form.triggerConfig.triggerType === 'command_received' && !form.triggerConfig.command) {
      missing.push('Command')
    }
  }

  if (form.actionService === 'Gmail') {
    if (!form.actionConfig.toEmail) missing.push('Email Address')
    if (!form.actionConfig.subject) missing.push('Email Subject')
  }

  if (form.actionService === 'Discord') {
    if (!form.actionConfig.webhookUrl) missing.push('Discord Webhook URL')
    if (!form.actionConfig.message) missing.push('Discord Message')
  }

  if (form.actionService === 'Spotify') {
    const playlistId = (form.actionConfig.playlistId || '').toString().trim()
    const sheetId = (form.actionConfig.spreadsheetId || form.triggerConfig?.spreadsheetId || '').toString().trim()
    const range = (form.actionConfig.range || form.triggerConfig?.range || '').toString().trim()
    const urlColumn = (form.actionConfig.urlColumn || '').toString().trim()

    if (!playlistId) missing.push('Spotify Playlist ID')
    if (!sheetId) missing.push('Spreadsheet ID')
    if (!range) missing.push('Sheet Range')
    if (!urlColumn) missing.push('Spotify Column')
  }

  if (form.actionService === 'Telegram') {
    if (!form.actionConfig.chatId) missing.push('Telegram Chat ID')
    if (!form.actionConfig.message) missing.push('Telegram Message')
  }

  if (form.intermediateActionService === 'OpenAI') {
    if (!form.intermediateActionConfig.prompt?.trim()) missing.push('OpenAI Prompt')
  }

  return missing.join(', ')
}

const loadRepositories = async () => {
  if (repositories.value.length > 0) return

  isLoadingRepositories.value = true
  try {
    repositories.value = await githubService.getRepositories()
    if (repositories.value.length === 0) {
      error.value = 'No repositories found. Make sure you have GitHub repositories and your GitHub account is linked.'
    }
  } catch (err) {
    console.error('Failed to load repositories:', err)
    const errorMessage = err instanceof Error ? err.message : 'Failed to load GitHub repositories'
    if (errorMessage.indexOf('GitHub username not configured') !== -1 || errorMessage.indexOf('No GitHub account linked') !== -1) {
      error.value = 'Please link your GitHub account first in your profile settings.'
    } else {
      error.value = `Failed to load repositories: ${errorMessage}`
    }
  } finally {
    isLoadingRepositories.value = false
  }
}

const onRepositoryChange = () => {
  if (form.triggerConfig.repositoryId) {
    const repo = repositories.value.find(r => r.id === parseInt(form.triggerConfig.repositoryId))
    if (repo) {
      form.description = `Monitor ${repo.full_name} for GitHub events and send email notifications`
    }
  }
}

const onNotificationTypeChange = () => {
  if (!form.triggerConfig.notificationTypes || form.triggerConfig.notificationTypes.length === 0) {
    form.triggerConfig.notificationTypes = ['push']
  }
}

const getSelectedRepositoryName = () => {
  if (!form.triggerConfig.repositoryId) return ''
  const repo = repositories.value.find(r => r.id === parseInt(form.triggerConfig.repositoryId))
  return repo?.full_name || ''
}

const isLoading = ref(false)
const error = ref<string | null>(null)
const isSendingTest = ref(false)
const isTestingWeather = ref(false)
const weatherTestResult = ref<any>(null)

const canSendTestEmail = computed(() => {
  return form.actionConfig.toEmail &&
         form.actionConfig.subject &&
         form.actionConfig.body
})

const sendTestEmail = async () => {
  if (!canSendTestEmail.value) return

  isSendingTest.value = true
  error.value = null

  try {
    const testEmailData = {
      to: form.actionConfig.toEmail,
      subject: form.actionConfig.subject,
      body: form.actionConfig.body
    }

    const response = await fetch('http://localhost:8080/test/email', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testEmailData)
    })

    const result = await response.json()

    if (response.ok) {
      alert('✅ Test email sent successfully!')
    } else {
      throw new Error(result.error || 'Failed to send test email')
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to send test email'
    console.error('Error sending test email:', err)
    alert('❌ Failed to send test email: ' + (err instanceof Error ? err.message : 'Unknown error'))
  } finally {
    isSendingTest.value = false
  }
}

const testWeatherTrigger = async () => {
  if (!form.triggerConfig.city) {
    alert('Please enter a city name first')
    return
  }

  isTestingWeather.value = true
  weatherTestResult.value = null

  try {
    const response = await fetch(`${API_BASE_URL}/test/weather`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        triggerConfig: {
          city: form.triggerConfig.city,
          temperature: form.triggerConfig.temperature || 0,
          condition: form.triggerConfig.condition || '',
          operator: form.triggerConfig.operator || 'greater_than'
        }
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    weatherTestResult.value = result

    if (result.success) {
      console.log('Weather trigger test result:', result)
    } else {
      console.error('Weather trigger test failed:', result.error)
    }
  } catch (err) {
    console.error('Error testing weather trigger:', err)
    alert('Failed to test weather trigger: ' + (err instanceof Error ? err.message : 'Unknown error'))
  } finally {
    isTestingWeather.value = false
  }
}

const testGoogleSheets = async () => {
  if (!canTestGoogleSheets.value) {
    alert('Please fill in the Spreadsheet ID and range first')
    return
  }

  isTestingGoogleSheets.value = true
  sheetsTestError.value = null

  try {
    const result = await areaService.testGoogleSheets({
      spreadsheetId: (form.triggerConfig?.spreadsheetId || '').toString().trim(),
      range: (form.triggerConfig?.range || '').toString().trim(),
    })

    sheetsTestResult.value = result
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Failed to test Google Sheets'
    sheetsTestError.value = message
    sheetsTestResult.value = null
    console.error('Error testing Google Sheets trigger:', err)
  } finally {
    isTestingGoogleSheets.value = false
  }
}

const formatDateTimeWithTimezone = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  const timezoneOffset = -date.getTimezoneOffset()
  const offsetHours = Math.floor(Math.abs(timezoneOffset) / 60)
  const offsetMinutes = Math.abs(timezoneOffset) % 60
  const offsetSign = timezoneOffset >= 0 ? '+' : '-'
  const offsetString = `${offsetSign}${String(offsetHours).padStart(2, '0')}:${String(offsetMinutes).padStart(2, '0')}`

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${offsetString}`
}

const createArea = async () => {
  if (!isFormValid.value) return

  isLoading.value = true
  error.value = null

  try {
    if (form.triggerService === 'GitHub' && form.actionService === 'Gmail') {
      await githubService.createArea(
        parseInt(form.triggerConfig.repositoryId),
        form.actionConfig.toEmail,
        form.triggerConfig.notificationTypes
      )
    } else if (form.triggerService === 'Google Sheets') {
      const isSpotifyAction = form.actionService === 'Spotify'
      const sanitizedActionConfig = isSpotifyAction
        ? {
            playlistId: (form.actionConfig.playlistId || '').toString().trim(),
            spreadsheetId: (form.actionConfig.spreadsheetId || form.triggerConfig?.spreadsheetId || '').toString().trim(),
            range: (form.actionConfig.range || form.triggerConfig?.range || '').toString().trim(),
            urlColumn: (form.actionConfig.urlColumn || 'SpotifyLink').toString().trim() || 'SpotifyLink',
            hasHeader: typeof form.actionConfig.hasHeader === 'boolean'
              ? form.actionConfig.hasHeader
              : !!form.triggerConfig?.hasHeader
          }
        : form.actionConfig

      const areaData: any = {
        name: form.areaName,
        description: form.description,
        triggerService: form.triggerService!,
        triggerType: 'SpreadsheetChange',
        actionService: form.actionService!,
        actionType: form.actionService === 'Gmail'
          ? 'SendEmail'
          : isSpotifyAction
          ? 'UpdatePlaylist'
          : 'Action',
        triggerConfig: {
          spreadsheetId: form.triggerConfig.spreadsheetId,
          sheetName: form.triggerConfig.sheetName,
          range: form.triggerConfig.range,
          hasHeader: !!form.triggerConfig.hasHeader
        },
        actionConfig: sanitizedActionConfig
      }

      if (form.intermediateActionService === 'OpenAI') {
        areaData.intermediateActionService = 'OpenAI'
        areaData.intermediateActionType = 'GenerateText'
        areaData.intermediateActionConfig = {
          prompt: form.intermediateActionConfig.prompt,
          systemPrompt: form.intermediateActionConfig.systemPrompt || '',
          temperature: form.intermediateActionConfig.temperature || 0.7,
          maxTokens: form.intermediateActionConfig.maxTokens || 500
        }
      }

      await areaService.createArea(areaData)
    } else if (form.triggerService === 'Weather') {
      const areaData: any = {
        name: form.areaName,
        description: form.description,
        triggerService: form.triggerService!,
        triggerType: 'Webhook',
        actionService: form.actionService!,
        actionType: form.actionService === 'Gmail'
          ? 'SendEmail'
          : form.actionService === 'Spotify'
          ? 'UpdatePlaylist'
          : 'Action',
        triggerConfig: form.triggerConfig,
        actionConfig: form.actionService === 'Spotify'
          ? {
              playlistId: (form.actionConfig.playlistId || '').toString().trim(),
              spreadsheetId: (form.actionConfig.spreadsheetId || '').toString().trim(),
              range: (form.actionConfig.range || '').toString().trim(),
              urlColumn: (form.actionConfig.urlColumn || 'SpotifyLink').toString().trim() || 'SpotifyLink',
              hasHeader: typeof form.actionConfig.hasHeader === 'boolean'
                ? form.actionConfig.hasHeader
                : true
            }
          : form.actionConfig
      }

      if (form.intermediateActionService === 'OpenAI') {
        areaData.intermediateActionService = 'OpenAI'
        areaData.intermediateActionType = 'GenerateText'
        areaData.intermediateActionConfig = {
          prompt: form.intermediateActionConfig.prompt,
          systemPrompt: form.intermediateActionConfig.systemPrompt || '',
          temperature: form.intermediateActionConfig.temperature || 0.7,
          maxTokens: form.intermediateActionConfig.maxTokens || 500
        }
      }

      await areaService.createArea(areaData)
    } else if (form.triggerService === 'Telegram') {
      const triggerConfig: any = {
        chatId: form.triggerConfig.chatId
      }

      if (form.triggerConfig.triggerType === 'keyword_match') {
        triggerConfig.keyword = form.triggerConfig.keyword
      } else if (form.triggerConfig.triggerType === 'command_received') {
        triggerConfig.command = form.triggerConfig.command
      }

      const areaData: any = {
        name: form.areaName,
        description: form.description,
        triggerService: form.triggerService!,
        triggerType: form.triggerConfig.triggerType || 'message_received',
        actionService: form.actionService!,
        actionType: form.actionService === 'Gmail'
          ? 'SendEmail'
          : form.actionService === 'Spotify'
          ? 'UpdatePlaylist'
          : 'Action',
        triggerConfig: triggerConfig,
        actionConfig: form.actionService === 'Spotify'
          ? {
              playlistId: (form.actionConfig.playlistId || '').toString().trim(),
              spreadsheetId: (form.actionConfig.spreadsheetId || '').toString().trim(),
              range: (form.actionConfig.range || '').toString().trim(),
              urlColumn: (form.actionConfig.urlColumn || 'SpotifyLink').toString().trim() || 'SpotifyLink',
              hasHeader: typeof form.actionConfig.hasHeader === 'boolean'
                ? form.actionConfig.hasHeader
                : true
            }
          : form.actionConfig
      }

      if (form.intermediateActionService === 'OpenAI') {
        areaData.intermediateActionService = 'OpenAI'
        areaData.intermediateActionType = 'GenerateText'
        areaData.intermediateActionConfig = {
          prompt: form.intermediateActionConfig.prompt,
          systemPrompt: form.intermediateActionConfig.systemPrompt || '',
          temperature: form.intermediateActionConfig.temperature || 0.7,
          maxTokens: form.intermediateActionConfig.maxTokens || 500
        }
      }

      await areaService.createArea(areaData)
    } else {
      let triggerConfig = { ...form.triggerConfig }

      if (isCalendarService(form.triggerService) && form.triggerConfig.eventTime) {
        const eventDateTime = new Date(form.triggerConfig.eventTime)
        const formattedDateTime = formatDateTimeWithTimezone(eventDateTime)

        const [datePart] = formattedDateTime.split('T')

        triggerConfig = {
          eventDate: datePart,
          eventTime: formattedDateTime,
          eventTitle: form.triggerConfig.eventTitle || '',
          calendarId: form.triggerConfig.calendarId || 'primary'
        }
      }

      const areaData: any = {
        name: form.areaName,
        description: form.description,
        triggerService: form.triggerService!,
        triggerType: isCalendarService(form.triggerService)
          ? 'Event'
          : form.triggerService === 'Spotify'
          ? 'Playback'
          : 'Webhook',
        actionService: form.actionService!,
        actionType: form.actionService === 'Gmail'
          ? 'SendEmail'
          : form.actionService === 'Spotify'
          ? 'UpdatePlaylist'
          : 'Action',
        triggerConfig: triggerConfig,
        actionConfig: form.actionService === 'Spotify'
          ? {
              playlistId: (form.actionConfig.playlistId || '').toString().trim(),
              spreadsheetId: (form.actionConfig.spreadsheetId || '').toString().trim(),
              range: (form.actionConfig.range || '').toString().trim(),
              urlColumn: (form.actionConfig.urlColumn || 'SpotifyLink').toString().trim() || 'SpotifyLink',
              hasHeader: typeof form.actionConfig.hasHeader === 'boolean'
                ? form.actionConfig.hasHeader
                : true
            }
          : form.actionConfig
      }

      if (form.intermediateActionService === 'OpenAI') {
        areaData.intermediateActionService = 'OpenAI'
        areaData.intermediateActionType = 'GenerateText'
        areaData.intermediateActionConfig = {
          prompt: form.intermediateActionConfig.prompt,
          systemPrompt: form.intermediateActionConfig.systemPrompt || '',
          temperature: form.intermediateActionConfig.temperature || 0.7,
          maxTokens: form.intermediateActionConfig.maxTokens || 500
        }
      }

      await areaService.createArea(areaData)
    }
    emit('save')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to create area'
    console.error('Error creating area:', err)
  } finally {
    isLoading.value = false
  }
}

const emit = defineEmits<{ (e: 'close'): void; (e: 'save'): void }>()
</script>

<style scoped>

.create-area-container {
  position: relative;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  width: 100%;
  background: var(--gradient-bg-primary);
  border-radius: 24px;
  overflow: hidden;
}


.main-card {
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 24px;
  padding: 0;
  width: 100%;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  position: relative;
}

[data-theme="light"] .main-card {
  background: #e5e7eb;
  border: 3px solid rgba(0, 0, 0, 0.4);
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 2rem 2rem 1rem 2rem;
  border-bottom: 1px solid var(--color-border-primary);
}

[data-theme="light"] .card-header {
  border-bottom: 3px solid rgba(0, 0, 0, 0.3);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.info-button {
  background: rgba(59, 130, 246, 0.15);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  padding: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.info-button:hover {
  background: rgba(59, 130, 246, 0.25);
  border-color: rgba(59, 130, 246, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.close-button {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--color-border-primary);
  border-radius: 8px;
  padding: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

[data-theme="light"] .close-button {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.3);
}

[data-theme="light"] .close-button :deep(.v-icon) {
  color: #1a1a1a !important;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.2);
}

[data-theme="light"] .close-button:hover {
  background: #f9fafb;
  border-color: rgba(0, 0, 0, 0.5);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.icon-container {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--spotify-green), var(--spotify-dark-green));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(29, 185, 84, 0.3);
}

.header-icon {
  color: white !important;
}

.header-text {
  flex: 1;
}

.card-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.02em;
}

[data-theme="light"] .card-title {
  color: #8b5cf6;
  text-shadow: 0 2px 4px rgba(139, 92, 246, 0.2);
  font-weight: 800;
}

.card-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.card-content {
  padding: 2rem;
}

.form-section {
  margin-bottom: 2.5rem;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.label-icon {
  color: var(--color-accent-primary) !important;
}

.label-text {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.input-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  letter-spacing: 0.01em;
}

.modern-input,
.modern-textarea {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 1rem;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
  transition: all 0.2s ease;
  outline: none;
}

[data-theme="light"] .modern-input,
[data-theme="light"] .modern-textarea {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.15);
  color: #1a1a1a;
}

.modern-input:focus,
.modern-textarea:focus {
  border-color: var(--color-accent-primary);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

[data-theme="light"] .modern-input:focus,
[data-theme="light"] .modern-textarea:focus {
  background: #ffffff;
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.modern-input::placeholder,
.modern-textarea::placeholder {
  color: var(--color-text-secondary);
  opacity: 0.7;
}

.modern-textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.connection-flow {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.service-selector {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 1.5rem;
  transition: all 0.2s ease;
}

[data-theme="light"] .service-selector {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.service-selector:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

[data-theme="light"] .service-selector:hover {
  background: #ffffff;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.selector-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.selector-icon {
  width: 40px;
  height: 40px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.selector-info {
  flex: 1;
}

.selector-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.01em;
}

.selector-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.service-dropdown {
  margin-top: 1rem;
}

.modern-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 1rem;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
  transition: all 0.2s ease;
  outline: none;
  width: 100%;
}

[data-theme="light"] .modern-select {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.15);
  color: #1a1a1a;
}

.modern-select:focus {
  border-color: var(--color-accent-primary);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

[data-theme="light"] .modern-select:focus {
  background: #ffffff;
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.selected-service {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.service-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-icon {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.service-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.placeholder-text {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  opacity: 0.7;
  font-style: italic;
}

.service-item {
  padding: 0.75rem 1rem !important;
  border-radius: 8px;
  margin: 0.25rem 0;
  transition: all 0.2s ease;
}

.service-item:hover {
  background: rgba(59, 130, 246, 0.1) !important;
}

.service-item-title {
  font-size: 0.875rem !important;
  font-weight: 500 !important;
  color: var(--color-text-primary) !important;
}

.connection-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0.5rem 0;
}

.arrow-line {
  flex: 1;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--color-accent-primary), transparent);
  opacity: 0.5;
}

.arrow-icon {
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 1rem;
}

.card-actions {
  padding: 1.5rem 2rem;
  border-top: 1px solid var(--color-border-primary);
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.05);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-primary);
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-text-primary);
}

.action-btn.primary {
  background: var(--gradient-accent);
  color: white;
  box-shadow: var(--shadow-glow);
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}

.action-btn.primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.config-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.config-btn, .test-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.config-btn {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.config-btn:hover {
  background: rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

.test-btn {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.test-btn:hover {
  background: rgba(34, 197, 94, 0.15);
  transform: translateY(-1px);
}

.temperature-input-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.operator-select {
  flex: 0 0 auto;
  min-width: 120px;
}

.temperature-input {
  flex: 1;
  min-width: 80px;
}

.temperature-unit {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
}

.test-section {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border-primary);
}

.test-weather-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.test-weather-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

.test-weather-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.test-result {
  margin-top: 1rem;
  padding: 1rem;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--color-border-primary);
}

.test-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.test-status.success {
  color: #22c55e;
}

.test-status.info {
  color: #3b82f6;
}

.weather-data {
  margin-top: 0.5rem;
  color: var(--color-text-secondary);
  font-size: 0.75rem;
}

.sheets-config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.sheets-checkbox {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.85rem 1rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 10px;
  cursor: pointer;
  user-select: none;
  transition: border-color 0.2s ease;
}

.sheets-checkbox:hover {
  border-color: var(--color-border-focus);
}

.sheets-checkbox input {
  width: 18px;
  height: 18px;
  accent-color: #3b82f6;
}

.sheets-checkbox span {
  color: var(--color-text-secondary);
  font-size: 0.85rem;
  font-weight: 500;
}

.sheets-test-actions {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.test-google-sheets-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  background: rgba(14, 165, 233, 0.12);
  color: #0ea5e9;
  border: 1px solid rgba(14, 165, 233, 0.28);
  align-self: flex-start;
}

.test-google-sheets-btn:hover:not(:disabled) {
  background: rgba(14, 165, 233, 0.16);
  transform: translateY(-1px);
}

.test-google-sheets-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.sheets-test-result {
  padding: 1rem;
  border-radius: 12px;
  border: 1px solid var(--color-border-primary);
  background: var(--color-bg-secondary);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.sheets-test-summary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.sheets-test-preview table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
  color: var(--color-text-primary);
}

.sheets-test-preview td {
  border: 1px solid var(--color-border-primary);
  padding: 0.4rem 0.6rem;
  background: rgba(255, 255, 255, 0.04);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .create-area-container {
    padding: 1rem;
  }

  .main-card {
    border-radius: 16px;
  }

  .card-header,
  .card-content,
  .card-actions {
    padding: 1.5rem;
  }

  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }

  .connection-flow {
    gap: 1rem;
  }

  .service-selector {
    padding: 1rem;
  }

  .card-actions {
    flex-direction: column;
  }

  .action-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-height: 800px) {
  .main-card {
    max-height: 80vh;
    overflow-y: auto;
  }

  .card-header {
    padding: 1.5rem 2rem 1rem 2rem;
  }

  .card-content {
    padding: 1.5rem 2rem;
  }

  .card-actions {
    padding: 1rem 2rem;
  }
}

.main-card {
  animation: slideInUp 0.4s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.service-selection {
  width: 100%;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-top: 16px;
}

.service-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  background: rgba(26, 31, 46, 0.6);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(10px);
  text-align: center;
}

.service-card:hover {
  background: rgba(26, 31, 46, 0.8);
  border-color: var(--color-border-focus);
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.service-card.more-services {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.service-card.more-services:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
}

.service-card-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.05);
}

.service-card-icon .service-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.service-card-name {
  color: var(--color-text-primary);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.2;
}

.selected-service-display {
  margin-top: 16px;
}

.selected-service-card {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: rgba(26, 31, 46, 0.8);
  border: 2px solid var(--color-accent-primary);
  border-radius: var(--radius-xl);
  backdrop-filter: blur(20px);
  box-shadow: var(--shadow-glow);
}

.service-info {
  flex: 1;
  margin-left: 12px;
}

.service-name {
  color: var(--color-text-primary);
  font-size: 16px;
  font-weight: 600;
  display: block;
}

.service-type {
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 400;
  display: block;
  margin-top: 2px;
}

.change-service-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-normal);
}

.change-service-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: var(--color-text-primary);
}

@media (max-width: 768px) {
  .service-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }

  .service-card {
    padding: 12px 8px;
  }

  .service-card-icon {
    width: 32px;
    height: 32px;
    margin-bottom: 6px;
  }

  .service-card-icon .service-icon {
    width: 20px;
    height: 20px;
  }

  .service-card-name {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .service-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

:deep(.v-field__outline) {
  --v-field-border-opacity: 0.1;
}

:deep(.v-field--variant-outlined .v-field__outline) {
  color: var(--color-border-primary);
}

:deep(.v-field--focused .v-field__outline) {
  color: var(--color-accent-primary);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  color: #ef4444;
  font-size: 0.875rem;
  margin: 1rem 0;
}

:deep(.v-list) {
  background: var(--color-bg-card) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
}

:deep(.v-list-item) {
  color: var(--color-text-primary) !important;
}

:deep(.v-list-item:hover) {
  background: var(--color-hover-bg) !important;
}

:deep(.v-field__input) {
  color: var(--color-text-primary) !important;
}


.config-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  transition: all 0.2s ease;
}

[data-theme="light"] .config-section {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.config-section:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

[data-theme="light"] .config-section:hover {
  background: #ffffff;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.config-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.config-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.config-icon .service-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.config-info {
  flex: 1;
}

.config-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.01em;
}

.config-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.config-content {
  padding-left: 3.5rem;
}

:deep(.v-field__outline) {
  color: var(--color-border-primary) !important;
}

.input-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.25rem;
  display: block;
  opacity: 0.8;
  font-style: italic;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.checkbox-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.875rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.checkbox-item:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
}

.checkbox-item input[type="checkbox"] {
  margin-right: 0.75rem;
  transform: scale(1.1);
}

.checkbox-label {
  font-weight: 500;
  color: var(--color-text-primary);
  font-size: 0.875rem;
}

.checkbox-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-left: 1.5rem;
  opacity: 0.8;
  font-style: italic;
}

.preview-section {
  margin-top: 2rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid var(--color-border-primary);
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.preview-icon {
  color: var(--color-accent-primary);
}

.preview-title {
  font-weight: 600;
  color: var(--color-text-primary);
}

.email-preview {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid var(--color-border-primary);
}

.email-header {
  margin-bottom: 0.5rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
}

.email-body {
  color: var(--color-text-secondary);
  line-height: 1.5;
  white-space: pre-wrap;
}

:deep(.v-field--focused .v-field__outline) {
  color: var(--color-accent-primary) !important;
}

.action-btn.test-email-btn {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: 0.5rem;
}

.action-btn.test-email-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #059669, #047857);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.action-btn.test-email-btn:disabled {
  background: #6b7280;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

:deep(.v-select .v-field__input) {
  color: var(--color-text-primary) !important;
}

:deep(.v-select .v-field__outline) {
  color: var(--color-border-primary) !important;
}


.info-box {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.875rem;
  margin-top: 1rem;
}

.info-box.warning {
  background: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.3);
}

.help-box {
  margin-top: 1rem;
  padding: 1rem;
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 12px;
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.help-content {
  flex: 1;
}

.help-content strong {
  color: rgba(255, 255, 255, 0.9);
  display: block;
  margin-bottom: 0.5rem;
}

.help-content ol {
  margin: 0.5rem 0 0 1.25rem;
  padding: 0;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.875rem;
  line-height: 1.6;
}

.help-content li {
  margin-bottom: 0.25rem;
}

.help-content code {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.125rem 0.375rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  color: #22c55e;
  font-size: 0.85rem;
}

.helper-link {
  color: #3b82f6;
  text-decoration: underline;
  font-weight: 600;
  transition: color 0.2s ease;
}

.helper-link:hover {
  color: #60a5fa;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.checkbox-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

.checkbox-item input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--color-accent-primary);
  cursor: pointer;
}

.checkbox-label {
  color: var(--color-text-primary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
}

.github-preview {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid var(--color-border-primary);
}

.preview-item {
  margin-bottom: 0.75rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
}

.preview-item:last-child {
  margin-bottom: 0;
}

.preview-item strong {
  color: var(--color-text-primary);
  font-weight: 600;
}

</style>
