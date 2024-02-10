<template>
  <div
    class="flex flex-col justify-center lg:w-[720px] items-center w-full h-screen p-5 lg:pl-20">
    <ez-card class="w-full">
      <div
        class="px-5 md:px-20 py-10 lg:py-20 flex flex-col justify-center items-center w-full">
        <div class="text-display mt-5 mb-10 lg:mb-20">Noxite</div>
        <div class="text-headline lg:py-5">
          {{ t('login.title') }}
        </div>
        <ez-textfield class="my-4 w-full">
          <div slot="label">{{ t('login.email') }}</div>
        </ez-textfield>
        <ez-textfield type="password" class="mb-1 w-full">
          <div slot="label">{{ t('login.password') }}</div>
        </ez-textfield>
        <div class="mt-4 w-full flex justify-end">
          <ez-text-button class="mr-5">
            {{ t('login.forgot-password') }}
          </ez-text-button>
          <ez-filled-button class="login-button">
            {{ t('login.login-action') }}
          </ez-filled-button>
        </div>
      </div>
    </ez-card>
  </div>
</template>
<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
<style lang="postcss">
.login-button {
  @apply flex items-center rounded-shape-corner-full;
  box-shadow: 0px 0px 40px 0px #ffffff80;
}
</style>
