<template>
  <div class="landing-dark">
    <!-- Desktop sidebar -->
    <v-navigation-drawer class="sidebar-desktop text-white" color="#0d0d0d" elevation="0" permanent rail>
      <v-list class="text-white" density="comfortable" nav lines="false">
        <v-tooltip text="Home" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-home" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Search" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-magnify" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Create" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-plus-circle" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Library" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-book-open-variant" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Profile" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-account-circle" class="text-white" rounded to="/login"></v-list-item>
          </template>
        </v-tooltip>
      </v-list>
    </v-navigation-drawer>

    <div class="content">
    <v-container class="pt-6 pb-4">
      <div class="d-flex align-center justify-space-between">
        <v-avatar size="44" class="gradient-avatar elevation-2">
          <v-icon color="white">mdi-account</v-icon>
        </v-avatar>
        <div class="d-flex ga-3 chips-row">
          <v-chip color="success" variant="flat" class="px-6" rounded="pill">All</v-chip>
          <v-chip variant="outlined" class="px-6" rounded="pill">My AREAs</v-chip>
          <v-chip variant="outlined" class="px-6" rounded="pill">Popular</v-chip>
          <v-chip variant="outlined" class="px-6" rounded="pill">Create</v-chip>
        </div>
        <v-btn icon variant="text"><v-icon>mdi-test-tube</v-icon></v-btn>
      </div>
    </v-container>

    <v-container>
      <h2 class="section-title mb-4">Popular AREAs</h2>
      <div class="cards-grid">
        <div class="card-col">
          <v-sheet class="area-card gradient-red" rounded="xl">
            <v-icon size="64" color="white">mdi-email-outline</v-icon>
          </v-sheet>
          <div class="card-title">Gmail → Discord</div>
          <div class="card-subtitle">Auto notification</div>
          <div class="card-description">Send Discord message when you receive important emails</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-green" rounded="xl">
            <v-icon size="64" color="white">mdi-music-note</v-icon>
          </v-sheet>
          <div class="card-title">Spotify → Twitter</div>
          <div class="card-subtitle">Auto sharing</div>
          <div class="card-description">Automatically tweet your favorite tracks</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-indigo" rounded="xl">
            <v-icon size="64" color="white">mdi-github</v-icon>
          </v-sheet>
          <div class="card-title">GitHub → Slack</div>
          <div class="card-subtitle">Team alerts</div>
          <div class="card-description">Notify your Slack channel when issues are opened</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-crimson" rounded="xl">
            <v-icon size="64" color="white">mdi-youtube</v-icon>
          </v-sheet>
          <div class="card-title">YouTube → Telegram</div>
          <div class="card-subtitle">New video</div>
          <div class="card-description">Get a Telegram ping when a channel uploads</div>
        </div>
      </div>
    </v-container>

    <v-container class="mt-6">
      <h2 class="section-title mb-4">Recommended for you</h2>
      <div class="cards-grid">
        <div class="card-col">
          <v-sheet class="area-card gradient-blue" rounded="xl">
            <v-icon size="64" color="white">mdi-weather-partly-cloudy</v-icon>
          </v-sheet>
          <div class="card-title">Weather → Telegram</div>
          <div class="card-subtitle">Daily reminder</div>
          <div class="card-description">Get weather forecast every morning on Telegram</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-pink" rounded="xl">
            <v-icon size="64" color="white">mdi-camera</v-icon>
          </v-sheet>
          <div class="card-title">Instagram → Dropbox</div>
          <div class="card-subtitle">Auto backup</div>
          <div class="card-description">Automatically backup your Instagram stories</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-teal" rounded="xl">
            <v-icon size="64" color="white">mdi-twitch</v-icon>
          </v-sheet>
          <div class="card-title">Twitch → Discord</div>
          <div class="card-subtitle">Go live</div>
          <div class="card-description">Alert your server when your stream starts</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-orange" rounded="xl">
            <v-icon size="64" color="white">mdi-newspaper-variant-outline</v-icon>
          </v-sheet>
          <div class="card-title">RSS → Notion</div>
          <div class="card-subtitle">Save articles</div>
          <div class="card-description">Append new posts to your Notion reading list</div>
        </div>
      </div>
    </v-container>

    <div class="bottom-nav">
      <div class="nav-inner">
        <v-btn variant="text" class="nav-btn" stacked>
          <v-icon size="26">mdi-home</v-icon>
          <span>Home</span>
        </v-btn>
        <v-btn variant="text" class="nav-btn" stacked>
          <v-icon size="26">mdi-magnify</v-icon>
          <span>Search</span>
        </v-btn>
        <v-btn variant="text" class="nav-btn" stacked>
          <v-icon size="26">mdi-plus-circle</v-icon>
          <span>Create</span>
        </v-btn>
        <v-btn variant="text" class="nav-btn" stacked>
          <v-icon size="26">mdi-book-open-variant</v-icon>
          <span>Library</span>
        </v-btn>
        <v-btn variant="text" class="nav-btn" stacked>
          <v-icon size="26">mdi-account-circle</v-icon>
          <span>Profile</span>
        </v-btn>
      </div>
    </div>
    </div>
  </div>

</template>

<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<style scoped>
.landing-dark {
  background: #000;
  color: white;
  min-height: 100vh;
}
.content {
  padding-left: 0;
}
.section-title {
  font-weight: 800;
  font-size: 32px;
}
.chips-row :deep(.v-chip) {
  font-weight: 600;
}
.gradient-avatar {
  background: linear-gradient(135deg, #7c4dff, #18ffff);
}
.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 24px;
}
.card-col { max-width: 320px; }
.area-card {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 24px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.25);
  transform: translateY(0) scale(1);
  background-size: 130% 130%;
  background-position: 50% 50%;
  transition:
    transform .25s ease,
    box-shadow .25s ease,
    background-position .6s ease,
    filter .25s ease;
}
.area-card :deep(.v-icon) {
  transition: transform .25s ease, opacity .25s ease;
}
.area-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 12px 28px rgba(0,0,0,0.35);
  background-position: 80% 20%;
}
.area-card:hover :deep(.v-icon) {
  transform: translateY(-2px) scale(1.06);
}
.area-card:active {
  transform: translateY(-2px) scale(0.99);
}

/* Subtle entrance animation */
.cards-grid .card-col { animation: fadeUp .45s ease both; }
.cards-grid .card-col:nth-child(2) { animation-delay: .05s; }
.cards-grid .card-col:nth-child(3) { animation-delay: .1s; }
.cards-grid .card-col:nth-child(4) { animation-delay: .15s; }

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Respect reduced motion */
@media (prefers-reduced-motion: reduce) {
  .area-card,
  .area-card :deep(.v-icon),
  .cards-grid .card-col {
    transition: none !important;
    animation: none !important;
  }
}
.gradient-red { background: linear-gradient(135deg, #c0392b, #d35400); }
.gradient-green { background: linear-gradient(135deg, #16a085, #27ae60); }
.gradient-blue { background: linear-gradient(135deg, #2980b9, #2ecc71); }
.gradient-pink { background: linear-gradient(135deg, #8e44ad, #e74c3c); }
.gradient-indigo { background: linear-gradient(135deg, #303f9f, #512da8); }
.gradient-crimson { background: linear-gradient(135deg, #c2185b, #e53935); }
.gradient-teal { background: linear-gradient(135deg, #00796b, #26a69a); }
.gradient-orange { background: linear-gradient(135deg, #ef6c00, #f39c12); }
.card-title { margin-top: 12px; font-weight: 800; font-size: 20px; }
.card-subtitle { color: rgba(255,255,255,0.85); font-weight: 700; }
.card-description { color: rgba(255,255,255,0.7); }
.bottom-nav {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  bottom: 16px;
  width: min(920px, 92%);
  background: rgba(255,255,255,0.06);
  border: 1px solid rgba(255,255,255,0.1);
  backdrop-filter: blur(12px);
  border-radius: 28px;
  padding: 8px 6px;
}
.nav-inner { display: grid; grid-template-columns: repeat(5, 1fr); gap: 6px; }
.nav-btn { color: white !important; text-transform: none; }
/* Sidebar visible on desktop, hidden on mobile */
.sidebar-desktop {
  display: none;
}
@media (min-width: 1280px) {
  .sidebar-desktop {
    display: block;
  }
  .content {
    padding-left: 80px; /* rail width */
  }
  .bottom-nav {
    display: none;
  }
}
</style>


