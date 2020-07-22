/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package iam

import (
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

func genIamResource(act ActionID, rscType TypeID, a *meta.ResourceAttribute) ([]types.Resource, error) {

	switch a.Basic.Type {
	case meta.Business:
		return genBusinessResource(act, rscType, a)
	case meta.DynamicGrouping:
		return genDynamicGroupingResource(act, rscType, a)
	case meta.EventPushing:
		return genEventSubscribeResource(act, rscType, a)
	case meta.EventWatch:
		return genResourceWatch(act, rscType, a)
	case meta.ProcessServiceTemplate, meta.ProcessTemplate:
		return genServiceTemplateResource(act, rscType, a)
	case meta.SetTemplate:
		return genSetTemplateResource(act, rscType, a)
	case meta.OperationStatistic:
		return genOperationStatisticResource(act, rscType, a)
	case meta.AuditLog:
		return genAuditLogResource(act, rscType, a)
	case meta.CloudAreaInstance:
		return genPlat(act, rscType, a)
	case meta.HostApply:
		return genHostApplyResource(act, rscType, a)
	case meta.CloudAccount:
		return genCloudAccountResource(act, rscType, a)
	case meta.CloudResourceTask:
		return genCloudResourceTaskResource(act, rscType, a)
	case meta.ResourcePoolDirectory:
		return genResourcePoolDirectoryResource(act, rscType, a)
	case meta.ProcessServiceInstance, meta.Process:
		return genProcessServiceInstanceResource(act, rscType, a)
	case meta.ModelModule, meta.ModelSet, meta.MainlineInstance, meta.MainlineInstanceTopology:
		return genBusinessTopologyResource(act, rscType, a)
	case meta.Model, meta.ModelAssociation:
		return genModelResource(act, rscType, a)
	case meta.ModelUnique:
		return genModelRelatedResource(act, rscType, a)
	case meta.ModelAttributeGroup:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelRelatedResource(act, rscType, a)
		}
	case meta.ModelClassification:
		return genModelClassificationResource(act, rscType, a)
	case meta.ModelInstance, meta.ModelInstanceAssociation:
		return genModelInstanceResource(act, rscType, a)
	case meta.AssociationType:
		return genAssociationTypeResource(act, rscType, a)
	case meta.ModelAttribute:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelAttributeResource(act, rscType, a)
		}
	case meta.ModelInstanceTopology, meta.MainlineModelTopology, meta.UserCustom:
		return genSkipResource(act, rscType, a)
	case meta.ConfigAdmin:
		return genGlobalConfigResource(act, rscType, a)
	case meta.MainlineModel:
		return genBusinessLayerResource(act, rscType, a)
	case meta.ModelTopology:
		return genModelTopologyViewResource(act, rscType, a)
	case meta.HostInstance:
		return genHostInstanceResource(act, rscType, a)

		// case meta.HostFavorite:
		// 	return genHostFavoriteResource(act, rscType, a)

		// case meta.SystemBase:
		// 	return new(types.Resource), nil
	case meta.ProcessServiceCategory:
		return genProcessServiceCategoryResource(act, rscType, a)
	}

	return nil, fmt.Errorf("gen id failed: unsupported resource type: %s", a.Type)
}

// generate business related resource id.
func genBusinessResource(act ActionID, typ TypeID, attribute *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create business do not related to instance authorize
	if act == CreateBusiness {
		return []types.Resource{r}, nil
	}

	// we have fuzzy authorize for frontend use, so we can not check this.
	// if attribute.InstanceID <= 0 {
	// 	return nil, errors.New("business instance id is 0")
	// }

	r.ID = strconv.FormatInt(attribute.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genDynamicGroupingResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to instance authorize
	if act == CreateBusinessCustomQuery || act == FindBusinessCustomQuery {
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	r.ID = att.InstanceIDEx

	// authorize based on business
	r.Attribute = map[string]interface{}{
		types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", Business, att.BusinessID)},
	}

	return []types.Resource{r}, nil
}

func genProcessServiceCategoryResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to instance authorize
	r.Type = types.ResourceType(Business)
	r.ID = strconv.FormatInt(att.BusinessID, 10)

	return []types.Resource{r}, nil
}

func genEventSubscribeResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateEventPushing {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genResourceWatch(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genServiceTemplateResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessServiceTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genSetTemplateResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessSetTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genOperationStatisticResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genAuditLogResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genPlat(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudArea {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types.Resource{r}, nil

}

func genHostApplyResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	r.Type = types.ResourceType(Business)
	r.ID = strconv.FormatInt(att.BusinessID, 10)

	return []types.Resource{r}, nil
}

func genCloudAccountResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudAccount {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types.Resource{r}, nil
}

func genCloudResourceTaskResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudResourceTask {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types.Resource{r}, nil
}

func genResourcePoolDirectoryResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateResourcePoolDirectory {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types.Resource{r}, nil
}

func genProcessServiceInstanceResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to exact service instance authorize
	r.Type = types.ResourceType(Business)
	r.ID = strconv.FormatInt(att.BusinessID, 10)

	return []types.Resource{r}, nil
}

func genBusinessTopologyResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to exact instance authorize
	r.Type = types.ResourceType(Business)
	r.ID = strconv.FormatInt(att.BusinessID, 10)

	return []types.Resource{r}, nil
}

func genModelResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to instance authorize
	if act == CreateSysModel {
		// create model authorized based on it's model group
		if len(att.Layers) > 0 {
			r.Type = types.ResourceType(SysModelGroup)
			r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
			return []types.Resource{r}, nil
		}
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genModelRelatedResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.Type = types.ResourceType(SysModel)
	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
	return []types.Resource{r}, nil

}

func genModelClassificationResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create model group do not related to instance authorize
	if act == CreateModelGroup {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genAssociationTypeResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateAssociationType {
		return []types.Resource{r}, nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)

	return []types.Resource{r}, nil
}

func genModelAttributeResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(EditSysModel),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)

	return []types.Resource{r}, nil
}

func genSkipResource(_ ActionID, _ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(Skip),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genGlobalConfigResource(_ ActionID, _ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(""),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genBusinessLayerResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genModelTopologyViewResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}
	return []types.Resource{r}, nil
}

func genHostInstanceResource(act ActionID, typ TypeID, a *meta.ResourceAttribute) ([]types.Resource, error) {

	// find host instances
	if act == Skip {
		r := types.Resource{
			System:    SystemIDCMDB,
			Type:      types.ResourceType(typ),
			Attribute: nil,
		}
		return []types.Resource{r}, nil
	}

	// transfer resource pool's host to it's another directory.
	if act == ResourcePoolHostTransferToDirectory {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysResourcePoolDirectory),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host in resource pool to business
	if act == ResourcePoolHostTransferToBusiness {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil

	}

	return []types.Resource{}, nil
}

func genBizModelAttributeResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(EditBusinessCustomField),
		Attribute: nil,
	}

	// if len(att.Layers) == 0 {
	// 	return nil, NotEnoughLayer
	// }

	r.ID = strconv.FormatInt(att.BusinessID, 10)

	return []types.Resource{r}, nil
}

func genModelInstanceResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// because we have to compatible to the any verify ,so we check the layers.
	// and if layer is 0, the exact authorize status is false as expected.
	if len(att.Layers) > 0 {
		if act == CreateSysInstance {
			r.Type = types.ResourceType(SysInstanceModel)
			r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
			return []types.Resource{r}, nil
		}

		r.ID = strconv.FormatInt(att.InstanceID, 10)

		// authorize based on a model
		r.Attribute = map[string]interface{}{
			types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", SysInstanceModel, att.Layers[0].InstanceID)},
		}
	}

	return []types.Resource{r}, nil
}