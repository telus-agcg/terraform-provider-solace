package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MsgVpnQueue struct for MsgVpnQueue
type MsgVpnQueue struct {
	AccessType                              *string         `tfsdk:"access_type"`
	ConsumerAckPropagationEnabled           *bool           `tfsdk:"consumer_ack_propagation_enabled"`
	DeadMsgQueue                            *string         `tfsdk:"dead_msg_queue"`
	DeliveryCountEnabled                    *bool           `tfsdk:"delivery_count_enabled"`
	DeliveryDelay                           *int64          `tfsdk:"delivery_delay"`
	EgressEnabled                           *bool           `tfsdk:"egress_enabled"`
	EventBindCountThreshold                 *EventThreshold `tfsdk:"event_bind_count_threshold"`
	EventMsgSpoolUsageThreshold             *EventThreshold `tfsdk:"event_msg_spool_usage_threshold"`
	EventRejectLowPriorityMsgLimitThreshold *EventThreshold `tfsdk:"event_reject_low_priority_msg_limit_threshold"`
	IngressEnabled                          *bool           `tfsdk:"ingress_enabled"`
	MaxBindCount                            *int64          `tfsdk:"max_bind_count"`
	MaxDeliveredUnackedMsgsPerFlow          *int64          `tfsdk:"max_delivered_unacked_msgs_per_flow"`
	MaxMsgSize                              *int32          `tfsdk:"max_msg_size"`
	MaxMsgSpoolUsage                        *int64          `tfsdk:"max_msg_spool_usage"`
	MaxRedeliveryCount                      *int64          `tfsdk:"max_redelivery_count"`
	MaxTtl                                  *int64          `tfsdk:"max_ttl"`
	MsgVpnName                              *string         `tfsdk:"msg_vpn_name"`
	Owner                                   *string         `tfsdk:"owner"`
	Permission                              *string         `tfsdk:"permission"`
	QueueName                               *string         `tfsdk:"queue_name"`
	RedeliveryEnabled                       *bool           `tfsdk:"redelivery_enabled"`
	RejectLowPriorityMsgEnabled             *bool           `tfsdk:"reject_low_priority_msg_enabled"`
	RejectLowPriorityMsgLimit               *int64          `tfsdk:"reject_low_priority_msg_limit"`
	RejectMsgToSenderOnDiscardBehavior      *string         `tfsdk:"reject_msg_to_sender_on_discard_behavior"`
	RespectMsgPriorityEnabled               *bool           `tfsdk:"respect_msg_priority_enabled"`
	RespectTtlEnabled                       *bool           `tfsdk:"respect_ttl_enabled"`
}

func (tfData *MsgVpnQueue) ToTF(apiData *sempv2.MsgVpnQueue) {
	AssignIfDstNotNil(&tfData.AccessType, apiData.AccessType)
	AssignIfDstNotNil(&tfData.ConsumerAckPropagationEnabled, apiData.ConsumerAckPropagationEnabled)
	AssignIfDstNotNil(&tfData.DeadMsgQueue, apiData.DeadMsgQueue)
	AssignIfDstNotNil(&tfData.DeliveryCountEnabled, apiData.DeliveryCountEnabled)
	AssignIfDstNotNil(&tfData.DeliveryDelay, apiData.DeliveryDelay)
	AssignIfDstNotNil(&tfData.EgressEnabled, apiData.EgressEnabled)
	AssignIfDstNotNil(&tfData.EventBindCountThreshold, EventThresholdToTF(apiData.EventBindCountThreshold))
	AssignIfDstNotNil(&tfData.EventMsgSpoolUsageThreshold, EventThresholdToTF(apiData.EventMsgSpoolUsageThreshold))
	AssignIfDstNotNil(&tfData.EventRejectLowPriorityMsgLimitThreshold, EventThresholdToTF(apiData.EventRejectLowPriorityMsgLimitThreshold))
	AssignIfDstNotNil(&tfData.IngressEnabled, apiData.IngressEnabled)
	AssignIfDstNotNil(&tfData.MaxBindCount, apiData.MaxBindCount)
	AssignIfDstNotNil(&tfData.MaxDeliveredUnackedMsgsPerFlow, apiData.MaxDeliveredUnackedMsgsPerFlow)
	AssignIfDstNotNil(&tfData.MaxMsgSize, apiData.MaxMsgSize)
	AssignIfDstNotNil(&tfData.MaxMsgSpoolUsage, apiData.MaxMsgSpoolUsage)
	AssignIfDstNotNil(&tfData.MaxRedeliveryCount, apiData.MaxRedeliveryCount)
	AssignIfDstNotNil(&tfData.MaxTtl, apiData.MaxTtl)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.Owner, apiData.Owner)
	AssignIfDstNotNil(&tfData.Permission, apiData.Permission)
	AssignIfDstNotNil(&tfData.QueueName, apiData.QueueName)
	AssignIfDstNotNil(&tfData.RedeliveryEnabled, apiData.RedeliveryEnabled)
	AssignIfDstNotNil(&tfData.RejectLowPriorityMsgEnabled, apiData.RejectLowPriorityMsgEnabled)
	AssignIfDstNotNil(&tfData.RejectLowPriorityMsgLimit, apiData.RejectLowPriorityMsgLimit)
	AssignIfDstNotNil(&tfData.RejectMsgToSenderOnDiscardBehavior, apiData.RejectMsgToSenderOnDiscardBehavior)
	AssignIfDstNotNil(&tfData.RespectMsgPriorityEnabled, apiData.RespectMsgPriorityEnabled)
	AssignIfDstNotNil(&tfData.RespectTtlEnabled, apiData.RespectTtlEnabled)
}

func (tfData *MsgVpnQueue) ToApi() *sempv2.MsgVpnQueue {
	return &sempv2.MsgVpnQueue{
		AccessType:                              tfData.AccessType,
		ConsumerAckPropagationEnabled:           tfData.ConsumerAckPropagationEnabled,
		DeadMsgQueue:                            tfData.DeadMsgQueue,
		DeliveryCountEnabled:                    tfData.DeliveryCountEnabled,
		DeliveryDelay:                           tfData.DeliveryDelay,
		EgressEnabled:                           tfData.EgressEnabled,
		EventBindCountThreshold:                 tfData.EventBindCountThreshold.ToApi(),
		EventMsgSpoolUsageThreshold:             tfData.EventMsgSpoolUsageThreshold.ToApi(),
		EventRejectLowPriorityMsgLimitThreshold: tfData.EventRejectLowPriorityMsgLimitThreshold.ToApi(),
		IngressEnabled:                          tfData.IngressEnabled,
		MaxBindCount:                            tfData.MaxBindCount,
		MaxDeliveredUnackedMsgsPerFlow:          tfData.MaxDeliveredUnackedMsgsPerFlow,
		MaxMsgSize:                              tfData.MaxMsgSize,
		MaxMsgSpoolUsage:                        tfData.MaxMsgSpoolUsage,
		MaxRedeliveryCount:                      tfData.MaxRedeliveryCount,
		MaxTtl:                                  tfData.MaxTtl,
		MsgVpnName:                              tfData.MsgVpnName,
		Owner:                                   tfData.Owner,
		Permission:                              tfData.Permission,
		QueueName:                               tfData.QueueName,
		RedeliveryEnabled:                       tfData.RedeliveryEnabled,
		RejectLowPriorityMsgEnabled:             tfData.RejectLowPriorityMsgEnabled,
		RejectLowPriorityMsgLimit:               tfData.RejectLowPriorityMsgLimit,
		RejectMsgToSenderOnDiscardBehavior:      tfData.RejectMsgToSenderOnDiscardBehavior,
		RespectMsgPriorityEnabled:               tfData.RespectMsgPriorityEnabled,
		RespectTtlEnabled:                       tfData.RespectTtlEnabled,
	}
}

// Terraform schema for MsgVpnQueue
func MsgVpnQueueSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnQueue",
		Attributes: map[string]tfsdk.Attribute{
			"access_type": {
				Type:        types.StringType,
				Description: "The access type for delivering messages to consumer flows bound to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"exclusive\"`. The allowed values and their meaning are:  <pre> \"exclusive\" - Exclusive delivery of messages to the first bound consumer flow. \"non-exclusive\" - Non-exclusive delivery of messages to all bound consumer flows in a round-robin fashion. </pre> ",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("exclusive", "non-exclusive"),
				},
			},
			"consumer_ack_propagation_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the propagation of consumer acknowledgements (ACKs) received on the active replication Message VPN to the standby replication Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"dead_msg_queue": {
				Type:        types.StringType,
				Description: "The name of the Dead Message Queue (DMQ) used by the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#DEAD_MSG_QUEUE\"`. Available since 2.2.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"delivery_count_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the ability for client applications to query the message delivery count of messages received from the Queue. This is a controlled availability feature. Please contact support to find out if this feature is supported for your use case. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.19.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"delivery_delay": {
				Type:        types.Int64Type,
				Description: "The delay, in seconds, to apply to messages arriving on the Queue before the messages are eligible for delivery. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.22.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"egress_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the transmission of messages from the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"event_bind_count_threshold": {
				Type:        EventThresholdType,
				Description: "",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"event_msg_spool_usage_threshold": {
				Type:        EventThresholdType,
				Description: "",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"event_reject_low_priority_msg_limit_threshold": {
				Type:        EventThresholdType,
				Description: "",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"ingress_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the reception of messages to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_bind_count": {
				Type:        types.Int64Type,
				Description: "The maximum number of consumer flows that can bind to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_delivered_unacked_msgs_per_flow": {
				Type:        types.Int64Type,
				Description: "The maximum number of messages delivered but not acknowledged per flow for the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_msg_size": {
				Type:        types.Int64Type,
				Description: "The maximum message size allowed in the Queue, in bytes (B). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000000`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_msg_spool_usage": {
				Type:        types.Int64Type,
				Description: "The maximum message spool usage allowed by the Queue, in megabytes (MB). A value of 0 only allows spooling of the last message received and disables quota checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `5000`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_redelivery_count": {
				Type:        types.Int64Type,
				Description: "The maximum number of times the Queue will attempt redelivery of a message prior to it being discarded or moved to the DMQ. A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"max_ttl": {
				Type:        types.Int64Type,
				Description: "The maximum time in seconds a message can stay in the Queue when `respectTtlEnabled` is `\"true\"`. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `maxTtl` configured for the Queue, is exceeded. A value of 0 disables expiry. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"owner": {
				Type:        types.StringType,
				Description: "The Client Username that owns the Queue and has permission equivalent to `\"delete\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"permission": {
				Type:        types.StringType,
				Description: "The permission level for all consumers of the Queue, excluding the owner. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"no-access\"`. The allowed values and their meaning are:  <pre> \"no-access\" - Disallows all access. \"read-only\" - Read-only access to the messages. \"consume\" - Consume (read and remove) messages. \"modify-topic\" - Consume messages or modify the topic/selector. \"delete\" - Consume messages, modify the topic/selector or delete the Client created endpoint altogether. </pre> ",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("no-access", "read-only", "consume", "modify-topic", "delete"),
				},
			},
			"queue_name": {
				Type:        types.StringType,
				Description: "The name of the Queue.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"redelivery_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable message redelivery. When enabled, the number of redelivery attempts is controlled by maxRedeliveryCount. When disabled, the message will never be delivered from the queue more than once. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"reject_low_priority_msg_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the checking of low priority messages against the `rejectLowPriorityMsgLimit`. This may only be enabled if `rejectMsgToSenderOnDiscardBehavior` does not have a value of `\"never\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"reject_low_priority_msg_limit": {
				Type:        types.Int64Type,
				Description: "The number of messages of any priority in the Queue above which low priority messages are not admitted but higher priority messages are allowed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"reject_msg_to_sender_on_discard_behavior": {
				Type:        types.StringType,
				Description: "Determines when to return negative acknowledgements (NACKs) to sending clients on message discards. Note that NACKs cause the message to not be delivered to any destination and Transacted Session commits to fail. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-queue-enabled\"`. The allowed values and their meaning are:  <pre> \"always\" - Always return a negative acknowledgment (NACK) to the sending client on message discard. \"when-queue-enabled\" - Only return a negative acknowledgment (NACK) to the sending client on message discard when the Queue is enabled. \"never\" - Never return a negative acknowledgment (NACK) to the sending client on message discard. </pre>  Available since 2.1.",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("always", "when-queue-enabled", "never"),
				},
			},
			"respect_msg_priority_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the respecting of message priority. When enabled, messages contained in the Queue are delivered in priority order, from 9 (highest) to 0 (lowest). MQTT queues do not support enabling message priority. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.8.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"respect_ttl_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the respecting of the time-to-live (TTL) for messages in the Queue. When enabled, expired messages are discarded or moved to the DMQ. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
