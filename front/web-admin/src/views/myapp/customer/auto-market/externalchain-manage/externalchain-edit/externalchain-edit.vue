<template>
  <create-external-chain-template></create-external-chain-template>
</template>
<script lang="ts">
import CreateExternalChainTemplate from "../../components/create-externalchain-template.vue";
export default {
  name: "",
  components: {
    CreateExternalChainTemplate
  },
  data() {
    return {};
  },
  created() {
    console.log("set --- created");
  },
  mounted() {
    console.log("set --- mounted");
  },
  methods: {
    editFn() {
      (this as any).$router.push({
        name:
          "myapp_customer_auto-market_externalchain-manage_externalchain-edit",
        query: {
          id: 2
        }
      });
    }
  }
};
</script>
<style lang="scss" scoped></style>
