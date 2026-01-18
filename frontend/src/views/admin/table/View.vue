<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">Table View</div>
    <div class="card-body overlay-wrapper">
      <div id="create_wrapper">
        <div class="tab-content" id="myTabContent">
          <div class="tab-pane fade show active" id="baseinfo">
            <div id="create_wrapper container">
              <div class="row">
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="ID" class="col-form-label col-sm-2 ">
                      ID:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.ID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="Name" class="col-form-label col-sm-2 ">
                      Name:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.Name }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="Content" class="col-form-label col-sm-2 ">
                      Content:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.Content }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="TableType" class="col-form-label col-sm-2 ">
                      TableType:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.TableType }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="ParentTable" class="col-form-label col-sm-2 ">
                      ParentTable:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.ParentTable }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="ParentField" class="col-form-label col-sm-2 ">
                      ParentField:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.ParentField }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="SelfField" class="col-form-label col-sm-2 ">
                      SelfField:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.SelfField }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="Sort" class="col-form-label col-sm-2 ">
                      Sort:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.Sort }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="Description" class="col-form-label col-sm-2 ">
                      Description:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.Description }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="Status" class="col-form-label col-sm-2 ">
                      Status:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.Status }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="CreatedBy" class="col-form-label col-sm-2 ">
                      CreatedBy:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.CreatedBy }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="UpdatedBy" class="col-form-label col-sm-2 ">
                      UpdatedBy:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.UpdatedBy }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="UpdatedAt" class="col-form-label col-sm-2 ">
                      UpdatedAt:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.UpdatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="CreatedAt" class="col-form-label col-sm-2 ">
                      CreatedAt:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.CreatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row mb-2">
                    <legend for="DeletedAt" class="col-form-label col-sm-2 ">
                      DeletedAt:
                    </legend>
                    <div class="col-sm-auto ">
                      {{ dataInfo.DeletedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3 mb-3">
                  <div class="form-group row">
                    <div class="col-sm-10">
                      <button class="btn btn-outline-secondary btn-sm" type="button" @click="goIndex">
                        Back
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getTable } from '@/api/table';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const dataInfo = ref({});

onMounted(() => {
  getTableData();
});

// get table list
const getTableData = async () => {
  const params = router.currentRoute.value.query;
  const res = await getTable({
    id: params.id,
  });
  dataInfo.value = res.data;
};

// goto index page
function goIndex() {
  router.push('/admin/table/index');
}
</script>

<style scoped></style>
